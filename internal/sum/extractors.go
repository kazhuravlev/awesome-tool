package sum

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

type URL struct{}

func (URL) Name() FactName {
	return "url"
}

func (URL) Deps() []FactName {
	return nil
}

func (URL) Extract(link *Link) {
	u, err := url.Parse(link.SrcLink.URL)
	if err != nil {
		return
	}

	link.Facts.Url.Data = u
}

type GitHub struct {
	// TODO: add github client, http client, credentials
}

func (GitHub) Name() FactName {
	return "github"
}

func (GitHub) Deps() []FactName {
	return []FactName{"url"}
}

func (GitHub) Extract(link *Link) {
	if !link.Facts.Url.HasData {
		return
	}

	u := link.Facts.Url.Data

	link.Facts.Url.Data = u
}

type Response struct {
	Client  *http.Client
	Timeout time.Duration
}

func (Response) Name() FactName {
	return "response"
}

func (Response) Deps() []FactName {
	return []FactName{"url"}
}

func (r *Response) Extract(link *Link) {
	ctx := context.Background()
	if r.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.Facts.Url.Data.String(), nil)
	if err != nil {
		return
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	duration := time.Since(start)

	link.Facts.Response = Fact[ResponseData]{
		HasData: true,
		Data: ResponseData{
			Protocol:   [2]int{resp.ProtoMajor, resp.ProtoMinor},
			Duration:   duration,
			StatusCode: resp.StatusCode,
			Body:       body,
			Headers:    resp.Header,
		},
	}
}
