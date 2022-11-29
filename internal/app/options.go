package app

import (
	"net/http"
	"time"

	"github.com/google/go-github/v48/github"
)

//go:generate options-gen -from-struct=Options -defaults-from=var

type Options struct {
	responseHttpClient *http.Client
	responseTimeout    time.Duration

	githubClient *github.Client
}

var defaultOptions = Options{
	responseHttpClient: http.DefaultClient,
	responseTimeout:    3 * time.Second,
}
