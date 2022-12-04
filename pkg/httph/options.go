package httph

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

//go:generate options-gen -from-struct=Options -defaults-from=var
type Options struct {
	client *http.Client

	rateLimitMapMu *sync.Mutex
	rateLimitMap   map[string]*rate.Limiter

	defaultRlConstructor func() *rate.Limiter
}

var defaultOptions = Options{}
