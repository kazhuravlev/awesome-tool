package source

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
)

func unmarshalFilename(filename string) (*Source, error) {
	sourceBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, errorsh.Wrap(err, "read source file")
	}

	var obj Source
	if err := yaml.Unmarshal(sourceBytes, &obj); err != nil {
		return nil, errorsh.Wrap(err, "unmarshal source file")
	}

	return &obj, nil
}
