package main

import (
	"fmt"

	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	sum.MustRegisterExtractor(sum.URL{})
	sum.MustRegisterExtractor(sum.GitHub{})

	const filename = "./examples/basic/data.yaml"
	sourceObj, err := source.ParseFile(filename)
	if err != nil {
		return errorsh.Wrap(err, "parse source file")
	}

	if err := source.Validate(*sourceObj); err != nil {
		return errorsh.Wrap(err, "validate source object")
	}

	for _, link := range sourceObj.Links {
		out, err := sum.GatherFacts(link)
		if err != nil {
			return err
		}

		fmt.Println(*out)

	}
	return nil
}
