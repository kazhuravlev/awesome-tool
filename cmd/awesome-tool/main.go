package main

import (
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/source"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	const filename = "./examples/basic/data.yaml"
	sourceObj, err := source.ParseFile(filename)
	if err != nil {
		return errorsh.Wrap(err, "parse source file")
	}

	if err := source.Validate(*sourceObj); err != nil {
		return errorsh.Wrap(err, "validate source object")
	}

	return nil
}
