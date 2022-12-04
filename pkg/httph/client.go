package httph

import (
	"net/http"
	"sync"

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

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

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

	return c.opts.client.Do(req)
}
