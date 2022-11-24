package sum

import (
	"net/url"

	"github.com/kazhuravlev/awesome-tool/internal/source"
)

const version = "1"

type FactName string

type Sum struct {
	Version string
	Rules   []source.Rule
	Groups  []source.Group
	Links   []Link
}

type Link struct {
	SrcLink source.Link
	Facts   LinkFacts
}

type LinkFacts struct {
	Url      Fact[*url.URL]
	Response Fact[struct{}]
	Github   Fact[struct{}]
	Gitlab   Fact[struct{}]
}

type Fact[T any] struct {
	HasData bool
	Data    T
}
