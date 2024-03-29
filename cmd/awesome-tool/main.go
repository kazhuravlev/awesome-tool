package main

import (
	"context"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/internal/app"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/pkg/httph"
	"github.com/urfave/cli/v3"
	"golang.org/x/oauth2"
	"log"
	"os"
	"path/filepath"
)

// Argument names
const (
	optSpecFilename      = "spec-file"
	optSumFilename       = "sum-file"
	optOutReadmeFilename = "out-readme"
	optGithubAccessToken = "github-token"
)

// Default argument values
const (
	optDefaultSpecFilename      = "./examples/basic/data.yaml"
	optDefaultSumFilename       = "./sum.yaml"
	optDefaultOutReadmeFilename = "./sum_readme.md"
)

func main() {
	app := &cli.App{ //nolint:exhaustruct
		Name: "awesome-tool",
		Commands: []*cli.Command{
			{
				Name:        "build",
				Description: "Build sum file from source",
				Action:      cmdBuild,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     optGithubAccessToken,
						Value:    "",
						Required: false,
					},
					&cli.StringFlag{
						Name:     optSpecFilename,
						Value:    optDefaultSpecFilename,
						Required: false,
					},
					&cli.StringFlag{
						Name:     optSumFilename,
						Value:    optDefaultSumFilename,
						Required: false,
					},
				},
			},
			{
				Name:        "render",
				Description: "Render sum file into template",
				Action:      cmdRender,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     optSumFilename,
						Value:    optDefaultSumFilename,
						Required: false,
					},
					&cli.StringFlag{
						Name:     optOutReadmeFilename,
						Value:    optDefaultOutReadmeFilename,
						Required: false,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func helpCreateApp(c *cli.Context) (*app.App, error) {
	githubAccessToken := c.String(optGithubAccessToken)

	httpClient, err := httph.New(httph.NewOptions())
	if err != nil {
		return nil, errorsh.Wrap(err, "create http instance")
	}

	var encoder app.Encoder
	switch ext := filepath.Ext(c.String(optSumFilename)); ext {
	default:
		return nil, errorsh.Newf("unknown out-sum filename extension: %s", ext)
	case ".yaml":
		encoder = app.YamlEncoder{}
	case ".json":
		encoder = app.JsonEncoder{}
	}

	github_ratelimit.NewRateLimitWaiterClient(oauth2.NewClient(
		context.TODO(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubAccessToken},
		),
	).Transport)
	appInst, err := app.New(app.NewOptions(
		app.WithGithubClient(
			github.NewClient(
				oauth2.NewClient(
					context.TODO(),
					oauth2.StaticTokenSource(
						&oauth2.Token{AccessToken: githubAccessToken},
					),
				),
			),
		),
		app.WithHttp(httpClient),
		app.WithMaxWorkers(10),
		app.WithSumEncoder(encoder),
	))
	if err != nil {
		return nil, errorsh.Wrap(err, "create app instance")
	}

	return appInst, nil
}

func cmdBuild(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInst, err := helpCreateApp(c)
	if err != nil {
		return errorsh.Wrap(err, "create application instance")
	}

	valSpecFilename := c.String(optSpecFilename)
	valSumFilename := c.String(optSumFilename)

	if err := appInst.Run(ctx, valSpecFilename, valSumFilename); err != nil {
		return errorsh.Wrap(err, "build sum")
	}

	return nil
}

func cmdRender(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInst, err := helpCreateApp(c)
	if err != nil {
		return errorsh.Wrap(err, "create application instance")
	}

	valSumFilename := c.String(optSumFilename)
	valOutReadmeFilename := c.String(optOutReadmeFilename)

	if err := appInst.Render(ctx, valSumFilename, valOutReadmeFilename); err != nil {
		return errorsh.Wrap(err, "render templates")
	}

	return nil
}
