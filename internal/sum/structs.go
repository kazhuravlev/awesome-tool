package sum

import (
	"net/url"
	"time"

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
	Response Fact[ResponseData]
	Github   Fact[struct{}]
}

type Fact[T any] struct {
	HasData bool
	Data    T
}

type ResponseData struct {
	// {2, 0} means http 2/0
	Protocol   [2]int
	Duration   time.Duration
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}
