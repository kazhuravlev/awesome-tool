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
	Url      *url.URL
	Response ResponseData
	Github   struct{}
}

type ResponseData struct {
	// {2, 0} means http 2/0
	Protocol   [2]int
	Duration   time.Duration
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}

type FactExtractor interface {
	// Name of extractor
	Name() FactName
	// Deps of this extractor
	Deps() []FactName
	// Implementation of extractor
	Extract(*Link) bool
}
