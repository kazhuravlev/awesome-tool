package app

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/kazhuravlev/awesome-tool/assets"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
	"github.com/kazhuravlev/just"
)

const outFilename = "sum.yaml"

type App struct {
	opts Options
}

func New(opts Options) (*App, error) {
	if err := opts.Validate(); err != nil {
		return nil, errorsh.Wrap(err, "bad configuration")
	}

	return &App{opts: opts}, nil
}

func (a *App) Run(ctx context.Context, filename string) error {
	sourceObj, err := source.ParseFile(filename)
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
			Client:  a.opts.responseHttpClient,
			Timeout: a.opts.responseTimeout,
		})
		facts.MustRegisterExtractor(&facts.GitHub{
			Client: a.opts.githubClient,
		})
	}

	linkFacts := make(map[int]facts.Facts, len(sourceObj.Links))
	for linkIdx := range sourceObj.Links {
		link := &sourceObj.Links[linkIdx]

		fmt.Printf("Gather facts about '%s'\n", link.URL)
		facts, err := facts.GatherFacts(ctx, *link)
		if err != nil {
			return errorsh.Wrapf(err, "gather facts for link '%s'", link.URL)
		}

		linkFacts[linkIdx] = *facts
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
					fmt.Println(link.Title, ":", rule.Name, check.Name(), ":", false, []string{"not all facts is collected"})
					// FIXME: add mark about dependency check
					linkChecks["__deps__"] = rules.CheckResult{
						CheckName: "Dependency check",
						IsPassed:  false,
						Errors:    []rules.Error{"Not all deps facts is collected"},
					}
					continue
				}

				checkResult := check.Test(link, linkFacts[linkIdx])
				fmt.Println(link.Title, ":", rule.Name, check.Name(), ":", checkResult)
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

		fmt.Println(string(sumYaml))

		if err := ioutil.WriteFile(outFilename, sumYaml, 0644); err != nil {
			return errorsh.Wrap(err, "write sum file")
		}
	}

	//fmt.Println(sumObj)
	// [ ] Apply rules
	// [ ] Render template + data

	return nil
}

func (a App) Render() error {
	var sumObj sum.Sum
	{
		bb, err := os.ReadFile(outFilename)
		if err != nil {
			return errorsh.Wrap(err, "read file")
		}

		if err := yaml.Unmarshal(bb, &sumObj); err != nil {
			return errorsh.Wrap(err, "unmarshal sum fil")
		}
	}

	bb, err := assets.FS.ReadFile("readme.md.tpl")
	if err != nil {
		return errorsh.Wrap(err, "read template")
	}

	reAnchor := regexp.MustCompile(`[^a-z0-9]+`)
	tmpl, err := template.New("readme.md").Funcs(template.FuncMap{
		"anchor": func(s string) string {
			return strings.Trim(reAnchor.ReplaceAllString(strings.ToLower(s), "-"), " -")
		},
		"add": func(n, x int) int {
			return n + x
		},
		"repeat": func(s string, n int) string {
			buf := bytes.NewBuffer(nil)
			for i := 0; i < n; i++ {
				buf.WriteString(s)
			}

			return buf.String()
		},
		"dict": func(values ...any) (map[string]any, error) {
			if len(values)%2 != 0 {
				return nil, errorsh.Newf("invalid dict call")
			}

			dict := make(map[string]any, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errorsh.Newf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}

			return dict, nil
		},
	}).Parse(string(bb))
	if err != nil {
		return errorsh.Wrap(err, "parse readme template")
	}

	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, sumObj); err != nil {
		return errorsh.Wrap(err, "exec template")
	}

	fmt.Println(buf.String())

	return nil
}
