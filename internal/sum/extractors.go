package sum

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/internal/source"
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

func (GitHub) Extract(ctx context.Context, link source.Link, facts *LinkFacts) (bool, error) {
	// FIXME: implement
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
