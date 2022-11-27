package app

import (
	"context"
	"fmt"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
	"net/http"
	"strconv"
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
			checkName, checkArgs, err := rules.ParseCheck(checkStringRaw)
			if err != nil {
				return errorsh.Wrapf(err, "parse check '%s'", checkStringRaw)
			}

			switch checkName {
			default:
				return errorsh.Newf("unknown check '%s'", checkName)
			case "response:status-eq":
				statusCode, err := strconv.Atoi(checkArgs)
				if err != nil {
					return errorsh.Wrapf(err, "bad arguments for check '%s'", checkName)
				}

				checks[checkStringRaw] = rules.CheckResponseStatusEq{
					ExpectedStatus: statusCode,
				}
			case "github-repo:stars-min":
				minimumStars, err := strconv.Atoi(checkArgs)
				if err != nil {
					return errorsh.Wrapf(err, "bad arguments for check '%s'", checkName)
				}

				checks[checkStringRaw] = rules.CheckResponseStatusEq{
					ExpectedStatus: minimumStars,
				}
			}
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

	fmt.Println(sumObj)
	// [ ] Apply rules
	// [ ] Render template + data

	return nil
}
