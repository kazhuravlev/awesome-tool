package app

import (
	"context"
	"fmt"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
	"net/http"
	"time"
)

func Run(ctx context.Context, filename string) error {
	// NOTE: register extractors
	{
		sum.MustRegisterExtractor(sum.URL{})
		sum.MustRegisterExtractor(&sum.Response{
			Client:  http.DefaultClient,
			Timeout: time.Second,
		})
		sum.MustRegisterExtractor(sum.GitHub{})
	}

	// NOTE: register checks
	{
		//rules.MustRegisterCheck(nil)
	}

	sourceObj, err := source.ParseFile(filename)
	if err != nil {
		return errorsh.Wrap(err, "parse source file")
	}

	if err := source.Validate(*sourceObj); err != nil {
		return errorsh.Wrap(err, "validate source object")
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
