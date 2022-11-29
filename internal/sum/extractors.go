package sum

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
)

const (
	FactUrl      FactName = "url"
	FactGithub   FactName = "github"
	FactResponse FactName = "response"
)

type URL struct{}

func (URL) Name() FactName {
	return FactUrl
}

func (URL) Deps() []FactName {
	return nil
}

func (URL) Extract(_ context.Context, link source.Link, facts *LinkFacts) (bool, error) {
	u, err := url.Parse(link.URL)
	if err != nil {
		return false, nil
	}

	facts.Url = u
	return true, nil
}

type GitHub struct {
	// TODO: add github client, http client, credentials
	Client *github.Client
}

func (GitHub) Name() FactName {
	return FactGithub
}

func (GitHub) Deps() []FactName {
	return []FactName{FactUrl}
}

func (e *GitHub) Extract(ctx context.Context, link source.Link, facts *LinkFacts) (bool, error) {
	if facts.Url.Host != "github.com" {
		return false, nil
	}

	parts := strings.Split(facts.Url.Path[1:], "/")
	if len(parts) != 2 {
		return false, errors.New("this is not a github repo url")
	}

	res, httpResp, err := e.Client.Repositories.Get(ctx, parts[0], parts[1])
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, errorsh.Wrap(err, "get repository details")
	}

	facts.Github = GithubData{
		OwnerUsername:    just.PointerUnwrap(res.Owner.Login),
		Name:             just.PointerUnwrap(res.Name),
		Description:      just.NullDefaultFalse(just.PointerUnwrapDefault(res.Description, "")),
		Homepage:         just.NullDefaultFalse(just.PointerUnwrapDefault(res.Homepage, "")),
		DefaultBranch:    just.PointerUnwrap(res.DefaultBranch),
		CreatedAt:        just.PointerUnwrap(res.CreatedAt).Time,
		PushedAt:         just.PointerUnwrap(res.PushedAt).Time,
		Language:         just.Null(just.PointerUnwrapDefault(res.Language, "")),
		Fork:             just.PointerUnwrap(res.Fork),
		ForksCount:       just.PointerUnwrap(res.ForksCount),
		NetworkCount:     just.PointerUnwrap(res.NetworkCount),
		OpenIssuesCount:  just.PointerUnwrap(res.OpenIssuesCount),
		StargazersCount:  just.PointerUnwrap(res.StargazersCount),
		SubscribersCount: just.PointerUnwrap(res.SubscribersCount),
		WatchersCount:    just.PointerUnwrap(res.WatchersCount),
		Topics:           res.Topics,
		Archived:         just.PointerUnwrap(res.Archived),
		Disabled:         just.PointerUnwrap(res.Disabled),
		License:          just.NullDefaultFalse(just.PointerUnwrapDefault(just.PointerUnwrapDefault(res.License, github.License{}).Name, "")),
	}
	return true, nil
}

type Response struct {
	Client  *http.Client
	Timeout time.Duration
}

func (Response) Name() FactName {
	return FactResponse
}

func (Response) Deps() []FactName {
	return []FactName{FactUrl}
}

func (r *Response) Extract(ctx context.Context, link source.Link, facts *LinkFacts) (bool, error) {
	if r.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, facts.Url.String(), nil)
	if err != nil {
		return false, nil
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	duration := time.Since(start)

	facts.Response = ResponseData{
		Protocol:   [2]int{resp.ProtoMajor, resp.ProtoMinor},
		Duration:   duration,
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}
	return true, nil
}
