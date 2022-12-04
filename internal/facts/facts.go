package facts

import (
	"context"
	"time"

	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
)

const version = "1"

type FactName string

type FactExtractor interface {
	// Name of extractor
	Name() FactName
	// Deps of this extractor
	Deps() []FactName
	// Extract implements of extractor
	Extract(context.Context, source.Link, *Data) (bool, error)
}

type Facts struct {
	Collected map[FactName]bool
	Data      Data
}

type Data struct {
	Url      UrlData
	Response ResponseData
	Github   GithubData
}

type UrlData struct {
	Scheme   string
	Hostname string
	Port     string
	Path     string
	Query    string
}

type ByteString []byte

func (bs ByteString) MarshalYAML() (any, error) {
	return string(bs), nil
}

func (bs *ByteString) UnmarshalYAML(b []byte) error {
	*bs = b
	return nil
}

type ResponseData struct {
	// {2, 0} means http 2/0
	Protocol        [2]int
	Duration        time.Duration
	StatusCode      int
	HtmlTitle       string
	HtmlDescription string
	Headers         map[string]string
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
