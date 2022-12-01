package app

import (
	"context"
	"fmt"

	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
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

	linkFacts := make([]facts.Facts, len(sourceObj.Links))
	for i := range sourceObj.Links {
		link := &sourceObj.Links[i]

		fmt.Printf("Gather facts about '%s'\n", link.URL)
		facts, err := facts.GatherFacts(ctx, *link)
		if err != nil {
			return errorsh.Wrapf(err, "gather facts for link '%s'", link.URL)
		}

		linkFacts[i] = *facts
	}

	rulesMap := make(map[source.RuleName]source.Rule)
	for _, rule := range sourceObj.Rules {
		rulesMap[rule.Name] = rule
	}

	// get global enabled rules
	globalRules := just.SliceMap(sourceObj.GlobalRulesEnabled, func(rn source.RuleName) source.Rule {
		return rulesMap[rn]
	})

	for linkIdx, link := range sourceObj.Links {
		// FIXME: implement group-level rules
		// FIXME: implement link-level rules
		linkRules := just.SliceCopy(globalRules)
		linkRules = just.SliceFilter(linkRules, func(rule source.Rule) bool {
			return !just.SliceContainsElem(link.RulesIgnored, rule.Name)
		})

		for _, rule := range linkRules {
			for _, checkStringRaw := range rule.Checks {
				check := checks[checkStringRaw]
				allFactsIsCollected := just.SliceAll(check.FactDeps(), func(factName facts.FactName) bool {
					return linkFacts[linkIdx].Collected[factName]
				})
				if !allFactsIsCollected {
					fmt.Println(link.Title, ":", rule.Name, check.Name(), ":", false, []string{"not all facts is collected"})
					continue
				}

				ok, errs := check.Test(link, linkFacts[linkIdx])
				fmt.Println(link.Title, ":", rule.Name, check.Name(), ":", ok, errs)
			}
		}
	}

	//fmt.Println(sumObj)
	// [ ] Apply rules
	// [ ] Render template + data

	return nil
}
