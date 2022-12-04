package main

import (
	"context"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/internal/app"
	"github.com/kazhuravlev/awesome-tool/pkg/httph"
	"golang.org/x/time/rate"
)

func main() {
	const filename = "./examples/basic/data.yaml"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient, err := httph.New(httph.NewOptions(
		httph.WithDefaultRlConstructor(func() *rate.Limiter {
			return rate.NewLimiter(rate.Every(time.Second), 5)
		}),
		httph.WithRateLimitMap(map[string]*rate.Limiter{
			"github.com": rate.NewLimiter(rate.Every(time.Second)/3, 2),
		}),
	))
	if err != nil {
		panic(err)
	}

	appInst, err := app.New(app.NewOptions(
		app.WithGithubClient(github.NewClient(nil)),
		app.WithHttp(httpClient),
		app.WithMaxWorkers(10),
	))
	if err != nil {
		panic(err)
	}

	if err := appInst.Run(ctx, filename); err != nil {
		panic(err)
	}
	// _ = ctx
	// if err := appInst.Render(); err != nil {
	// 	panic(err)
	// }
}
