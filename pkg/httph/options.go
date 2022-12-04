package httph

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

//go:generate options-gen -from-struct=Options -defaults-from=var
type Options struct {
	client *http.Client

	rateLimitMap map[string]*rate.Limiter

	defaultRlConstructor func() *rate.Limiter
}

var defaultOptions = Options{
	client:               http.DefaultClient,
	rateLimitMap:         make(map[string]*rate.Limiter),
	defaultRlConstructor: func() *rate.Limiter { return rate.NewLimiter(rate.Every(time.Second), 5) },
}
