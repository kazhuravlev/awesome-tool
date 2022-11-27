package app

import (
	"context"
	"fmt"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
	"net/http"
	"time"
)

func Run(ctx context.Context, filename string) error {
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
		sum.MustRegisterExtractor(sum.URL{})
		sum.MustRegisterExtractor(&sum.Response{
			Client:  http.DefaultClient,
			Timeout: 3 * time.Second,
		})
		sum.MustRegisterExtractor(sum.GitHub{})
	}

	sumObj, err := sum.GatherFacts(ctx, *sourceObj)
	if err != nil {
		return errorsh.Wrap(err, "gather facts for source obj")
	}

	for _, link := range sumObj.Links {
		for _, rule := range link.Rules {
			for _, checkStringRaw := range rule.Checks {
				check := checks[checkStringRaw]
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
