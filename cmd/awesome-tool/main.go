package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/internal/app"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/pkg/httph"
	"github.com/urfave/cli/v3"
	"golang.org/x/time/rate"
)

const inFilename = "./examples/basic/data.yaml"
const outFilename = "./sum.yaml"

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name: "awesome-tool",
		Commands: []*cli.Command{
			{
				Name:        "build",
				Description: "Build sum file from source",
				Action:      cmdBuild,
			},
			{
				Name:        "render",
				Description: "Render sum file into template",
				Action:      cmdRender,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func helpCreateApp() (*app.App, error) {
	httpClient, err := httph.New(httph.NewOptions(
		httph.WithDefaultRlConstructor(func() *rate.Limiter {
			return rate.NewLimiter(rate.Every(time.Second), 5)
		}),
		httph.WithRateLimitMap(map[string]*rate.Limiter{
			"github.com": rate.NewLimiter(rate.Every(time.Second)/1, 2),
		}),
	))
	if err != nil {
		return nil, errorsh.Wrap(err, "create http instance")
	}

	appInst, err := app.New(app.NewOptions(
		app.WithGithubClient(github.NewClient(nil)),
		app.WithHttp(httpClient),
		app.WithMaxWorkers(10),
	))
	if err != nil {
		return nil, errorsh.Wrap(err, "create app instance")
	}

	return appInst, nil
}

func cmdBuild(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInst, err := helpCreateApp()
	if err != nil {
		return errorsh.Wrap(err, "create application instance")
	}

	if err := appInst.Run(ctx, inFilename, outFilename); err != nil {
		return errorsh.Wrap(err, "build sum")
	}

	return nil
}

func cmdRender(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInst, err := helpCreateApp()
	if err != nil {
		return errorsh.Wrap(err, "create application instance")
	}

	if err := appInst.Render(ctx, outFilename); err != nil {
		return errorsh.Wrap(err, "render templates")
	}

	return nil
}
