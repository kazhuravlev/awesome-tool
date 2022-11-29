package main

import (
	"context"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/internal/app"
)

func main() {
	const filename = "./examples/basic/data.yaml"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInst, err := app.New(app.NewOptions(
		app.WithGithubClient(github.NewClient(nil)),
	))
	if err != nil {
		panic(err)
	}

	if err := appInst.Run(ctx, filename); err != nil {
		panic(err)
	}
}
