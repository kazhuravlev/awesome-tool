package app

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
)

func tplDict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errorsh.Newf("invalid dict call")
	}

	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errorsh.Newf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}

	return dict, nil
}

func tplRepeat(s string, n int) string {
	buf := bytes.NewBuffer(nil)
	for i := 0; i < n; i++ {
		buf.WriteString(s)
	}

	return buf.String()
}

func tplAdd(n, x int) int { return n + x }

func tplAnchor(s string) string {
	return strings.Trim(reAnchor.ReplaceAllString(strings.ToLower(s), "-"), " -")
}

var tplFuncLib = template.FuncMap{
	"anchor": tplAnchor,
	"add":    tplAdd,
	"repeat": tplRepeat,
	"dict":   tplDict,
}
