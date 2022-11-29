package sum

import (
	"context"
	"net/url"
	"time"

	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
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
	SrcLink        source.Link
	FactsCollected map[FactName]bool
	Facts          LinkFacts
	// NOTE: This is a duplicate for each link. It is a result set of rules,
	//   which will applied to exact this link after all enable/disable/ignore
	//   rules.
	Rules []source.Rule
}

type LinkFacts struct {
	Url      *url.URL
	Response ResponseData
	Github   GithubData
}

type ResponseData struct {
	// {2, 0} means http 2/0
	Protocol   [2]int
	Duration   time.Duration
	StatusCode int
	Body       []byte
	Headers    map[string][]string
}

type GithubData struct {
	OwnerUsername    string
	Name             string
	Description      just.NullVal[string]
	Homepage         just.NullVal[string]
	DefaultBranch    string
	CreatedAt        time.Time
	PushedAt         time.Time
	Language         just.NullVal[string]
	Fork             bool
	ForksCount       int
	NetworkCount     int
	OpenIssuesCount  int
	StargazersCount  int
	SubscribersCount int
	WatchersCount    int
	Topics           []string
	Archived         bool
	Disabled         bool
	License          just.NullVal[string]
}

type FactExtractor interface {
	// Name of extractor
	Name() FactName
	// Deps of this extractor
	Deps() []FactName
	// Extract implements of extractor
	Extract(context.Context, source.Link, *LinkFacts) (bool, error)
}
