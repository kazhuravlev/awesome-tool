package yamlh

import (
	"github.com/goccy/go-yaml"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"os"
)

// UnmarshalFile helps to read file and then unmarshal it into obj.
func UnmarshalFile(filename string, obj any) error {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return errorsh.Wrap(err, "read file")
	}

	if err := yaml.Unmarshal(bb, obj); err != nil {
		return errorsh.Wrap(err, "unmarshal")
	}

	return nil
}
