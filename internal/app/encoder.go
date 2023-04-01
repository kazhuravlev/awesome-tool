package app

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/pkg/yamlh"
	"os"
)

type YamlEncoder struct{}

func (e YamlEncoder) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (e YamlEncoder) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (e YamlEncoder) MarshalFile(filename string, v any) error {
	return yamlh.MarshalFile(filename, v)
}

func (e YamlEncoder) UnmarshalFile(filename string, v any) error {
	return yamlh.UnmarshalFile(filename, v)
}

type JsonEncoder struct{}

func (e JsonEncoder) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (e JsonEncoder) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (e JsonEncoder) UnmarshalFile(filename string, obj any) error {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return errorsh.Wrap(err, "read file")
	}

	if err := json.Unmarshal(bb, obj); err != nil {
		return errorsh.Wrap(err, "unmarshal")
	}

	return nil
}

func (e JsonEncoder) MarshalFile(filename string, obj any) error {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(obj); err != nil {
		return errorsh.Wrap(err, "marshal")
	}

	if err := os.WriteFile(filename, buf.Bytes(), 0o0644); err != nil {
		return errorsh.Wrap(err, "write file")
	}

	return nil
}
