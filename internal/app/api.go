package app

import (
	"context"
	"fmt"

	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
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
		sum.MustRegisterExtractor(&sum.URL{})
		sum.MustRegisterExtractor(&sum.Response{
			Client:  a.opts.responseHttpClient,
			Timeout: a.opts.responseTimeout,
		})
		sum.MustRegisterExtractor(&sum.GitHub{
			Client: a.opts.githubClient,
		})
	}

	sumObj, err := sum.GatherFacts(ctx, *sourceObj)
	if err != nil {
		return errorsh.Wrap(err, "gather facts for source obj")
	}

	for _, link := range sumObj.Links {
		for _, rule := range link.Rules {
			for _, checkStringRaw := range rule.Checks {
				check := checks[checkStringRaw]
				allFactsIsCollected := just.SliceAll(check.FactDeps(), func(factName sum.FactName) bool {
					return link.FactsCollected[factName]
				})
				if !allFactsIsCollected {
					fmt.Println(link.SrcLink.Title, ":", rule.Name, check.Name(), ":", false, []string{"not all facts is collected"})
					continue
				}

				ok, errs := check.Test(link)
				fmt.Println(link.SrcLink.Title, ":", rule.Name, check.Name(), ":", ok, errs)
			}
		}
	}

	//fmt.Println(sumObj)
	// [ ] Apply rules
	// [ ] Render template + data

	return nil
}
