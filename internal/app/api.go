package app

import (
	"context"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
	"github.com/kazhuravlev/just"
)

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
	linksChecks := make(map[int]map[string][]rules.Error, len(sourceObj.Links))
	for linkIdx, link := range sourceObj.Links {
		// FIXME: implement group-level rules
		// FIXME: implement link-level rules
		linkRules := just.SliceCopy(globalRules)
		linkRules = just.SliceFilter(linkRules, func(rule source.Rule) bool {
			return !just.SliceContainsElem(link.RulesIgnored, rule.Name)
		})

		linksRules[linkIdx] = linkRules

		linkChecks := make(map[string][]rules.Error, len(linkRules))
		for _, rule := range linkRules {
			for _, checkStringRaw := range rule.Checks {
				check := checks[checkStringRaw]
				allFactsIsCollected := just.SliceAll(check.FactDeps(), func(factName facts.FactName) bool {
					return linkFacts[linkIdx].Collected[factName]
				})
				if !allFactsIsCollected {
					fmt.Println(link.Title, ":", rule.Name, check.Name(), ":", false, []string{"not all facts is collected"})
					linkChecks["__deps__"] = []rules.Error{"Not all deps facts is collected"}
					continue
				}

				ok, errs := check.Test(link, linkFacts[linkIdx])
				fmt.Println(link.Title, ":", rule.Name, check.Name(), ":", ok, errs)
				linkChecks[checkStringRaw] = errs
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
	}

	//fmt.Println(sumObj)
	// [ ] Apply rules
	// [ ] Render template + data

	return nil
}
