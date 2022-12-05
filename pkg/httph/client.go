package httph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"golang.org/x/time/rate"
)

type Client struct {
	opts Options

	rateLimitMapMu *sync.Mutex
}

func New(opts Options) (*Client, error) {
	if err := opts.Validate(); err != nil {
		return nil, errorsh.Wrap(err, "bad configuration")
	}

	return &Client{
		opts:           opts,
		rateLimitMapMu: new(sync.Mutex),
	}, nil
}

var ctxKeyEquivRedirectNum = struct{}{}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	equivRedirectNum := 0
	if val, ok := ctx.Value(ctxKeyEquivRedirectNum).(int); ok {
		equivRedirectNum = val
	}

	hostname := req.URL.Hostname()
	var rl *rate.Limiter
	func() {
		c.rateLimitMapMu.Lock()
		defer c.rateLimitMapMu.Unlock()

		rateLimiter, ok := c.opts.rateLimitMap[hostname]
		if !ok {
			rl = c.opts.defaultRlConstructor()
			return
		}

		c.opts.rateLimitMap[hostname] = rateLimiter
		rl = rateLimiter
	}()

	if err := rl.Wait(ctx); err != nil {
		return nil, err
	}

	resp, err := c.opts.client.Do(req)
	if err != nil {
		return nil, errorsh.Wrap(err, "do original request")
	}

	buf := bytes.NewBuffer(nil)
	teeReader := io.TeeReader(resp.Body, buf)
	body, err := io.ReadAll(teeReader)
	if err != nil {
		return nil, errorsh.Wrap(err, "read response body")
	}

	resp.Body.Close()
	resp.Body = io.NopCloser(buf)

	if c.opts.maxEquivRedirects > 0 {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			// NOTE: we just unable to read response as html. It is ok, no error.if err != nil {

		} else {
			const noEquiv = "__no__equiv__"
			equiv := doc.Find("html>head>meta[http-equiv=refresh]").First().AttrOr("content", noEquiv)
			if equiv != noEquiv {
				groups := regexp.MustCompile(`url=(?P<url>.*)`).FindAllStringSubmatch(equiv, -1)
				if len(groups) == 1 && len(groups[0]) == 2 {
					redirectUrl := groups[0][1]
					u, err := url.Parse(redirectUrl)
					if err == nil {
						if equivRedirectNum >= c.opts.maxEquivRedirects {
							return nil, errorsh.Newf("max attempts to fetch url")
						}

						equivRedirectNum += 1

						fmt.Printf("Additional attempt based on http-equiv. Replace URL '%s' to '%s'\n", req.URL.String(), redirectUrl)
						req.URL = u

						ctx2 := context.WithValue(ctx, ctxKeyEquivRedirectNum, equivRedirectNum)

						return c.Do(req.WithContext(ctx2))
					}
				}
			}
		}
	}

	return resp, nil
}
