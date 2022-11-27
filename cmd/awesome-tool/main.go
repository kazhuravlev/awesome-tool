package main

import (
	"context"
	"github.com/kazhuravlev/awesome-tool/internal/app"
)

func main() {
	const filename = "./examples/basic/data.yaml"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Run(ctx, filename); err != nil {
		panic(err)
	}
}
