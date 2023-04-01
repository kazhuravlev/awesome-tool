package app

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kazhuravlev/awesome-tool/pkg/yamlh"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
	"sync/atomic"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/kazhuravlev/awesome-tool/assets"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
	"github.com/kazhuravlev/just"
	"golang.org/x/sync/semaphore"
)

var (
	tplReadme = just.Must(template.New("readme.md").Funcs(tplFuncLib).Parse(string(just.Must(assets.FS.ReadFile("readme.go.tpl")))))
)

var reAnchor = regexp.MustCompile(`[^a-z0-9]+`)

type App struct {
	opts Options
}

func New(opts Options) (*App, error) {
	if err := opts.Validate(); err != nil {
		return nil, errorsh.Wrap(err, "bad configuration")
	}

	return &App{opts: opts}, nil
}

func (a *App) Run(ctx context.Context, inFilename, outFilename string) error {
	sourceObj, err := source.ParseFile(inFilename)
	if err != nil {
		return errorsh.Wrap(err, "parse source file")
	}

	if err := source.Validate(*sourceObj); err != nil {
		return errorsh.Wrap(err, "validate source object")
	}

	checks := make(map[string]rules.Check)
	for _, rule := range sourceObj.Rules {
		for _, checkStringRaw := range rule.Checks {
			check, err := rules.ParseCheck(checkStringRaw)
			if err != nil {
				return errorsh.Wrapf(err, "parse check '%s'", checkStringRaw)
			}

			checks[checkStringRaw] = check
		}
	}

	// NOTE: register extractors
	{
		// TODO: extract uniq deps from current checks and use only required
		//   fact extractors.
		facts.MustRegisterExtractor(&facts.URL{})
		facts.MustRegisterExtractor(&facts.Response{
			Client:  a.opts.http,
			Timeout: a.opts.responseTimeout,
		})
		facts.MustRegisterExtractor(&facts.GitHub{
			Client: a.opts.githubClient,
		})
		facts.MustRegisterExtractor(&facts.Meta{})
	}

	linkFactsMu := new(sync.Mutex)
	linkFacts := make(map[int]facts.Facts, len(sourceObj.Links))
	sem := semaphore.NewWeighted(int64(a.opts.maxWorkers))
	var counter int64
	for linkIdx := range sourceObj.Links {
		link := &sourceObj.Links[linkIdx]

		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		go func(linkIdx int, link *source.Link) {
			defer sem.Release(1)

			curValue := atomic.AddInt64(&counter, 1)
			fmt.Printf("Gather facts [%d/%d] about '%s'\n", curValue, len(sourceObj.Links), link.URL)
			facts, err := facts.GatherFacts(ctx, *link)
			if err != nil {
				log.Printf("fail to gather facts: %s", err)
				return
			}

			linkFactsMu.Lock()
			linkFacts[linkIdx] = *facts
			linkFactsMu.Unlock()
		}(linkIdx, link)
	}
	if err := sem.Acquire(ctx, int64(a.opts.maxWorkers)); err != nil {
		return errorsh.Wrap(err, "wait to workers finished")
	}

	rulesMap := make(map[source.RuleName]source.Rule)
	for _, rule := range sourceObj.Rules {
		rulesMap[rule.Name] = rule
	}

	// get global enabled rules
	globalRules := just.SliceMap(sourceObj.GlobalRulesEnabled, func(rn source.RuleName) source.Rule {
		return rulesMap[rn]
	})

	linksRules := make(map[int][]source.Rule, len(sourceObj.Links))
	linksChecks := make(map[int]map[string]rules.CheckResult, len(sourceObj.Links))
	for linkIdx, link := range sourceObj.Links {
		// FIXME: implement group-level rules
		// FIXME: implement link-level rules
		linkRules := just.SliceCopy(globalRules)
		linkRules = just.SliceFilter(linkRules, func(rule source.Rule) bool {
			return !just.SliceContainsElem(link.RulesIgnored, rule.Name)
		})

		linksRules[linkIdx] = linkRules

		linkChecks := make(map[string]rules.CheckResult, len(linkRules))
		for _, rule := range linkRules {
			for _, checkStringRaw := range rule.Checks {
				check := checks[checkStringRaw]
				allFactsIsCollected := just.SliceAll(check.FactDeps(), func(factName facts.FactName) bool {
					return linkFacts[linkIdx].Collected[factName]
				})
				if !allFactsIsCollected {
					fmt.Println(link.URL, ":", rule.Name, check.Name(), ":", false, []string{"not all facts is collected"})
					// FIXME: add mark about dependency check
					linkChecks["__deps__"] = rules.CheckResult{
						CheckName: "Dependency check",
						IsPassed:  false,
						Errors:    []rules.Error{"Not all deps facts is collected"},
					}
					continue
				}

				checkResult := check.Test(link, linkFacts[linkIdx])
				fmt.Println(link.URL, ":", rule.Name, check.Name(), ":", checkResult)
				linkChecks[checkStringRaw] = checkResult
			}
		}

		linksChecks[linkIdx] = linkChecks
	}

	{
		sumObj := sum.Build(*sourceObj, linksRules, linkFacts, linksChecks)
		sumYaml, err := yaml.Marshal(sumObj)
		if err != nil {
			return errorsh.Wrap(err, "marshal sum object")
		}

		if err := ioutil.WriteFile(outFilename, sumYaml, 0644); err != nil {
			return errorsh.Wrap(err, "write sum file")
		}
	}

	return nil
}

func (a App) Render(ctx context.Context, sumFilename, readmeFilename string) error {
	var sumObj sum.Sum
	if err := yamlh.UnmarshalFile(sumFilename, &sumObj); err != nil {
		return errorsh.Wrap(err, "unmarshal sum file")
	}

	buf := bytes.NewBuffer(nil)
	if err := tplReadme.Execute(buf, sumObj); err != nil {
		return errorsh.Wrap(err, "exec template")
	}

	if err := os.WriteFile(readmeFilename, buf.Bytes(), 0644); err != nil {
		return errorsh.Wrap(err, "write result readme file")
	}

	return nil
}
