package app

import (
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/pkg/httph"
	"github.com/kazhuravlev/just"
)

//go:generate options-gen -from-struct=Options -defaults-from=var

type Options struct {
	responseTimeout time.Duration

	githubClient *github.Client

	http *httph.Client
}

var defaultOptions = Options{
	responseTimeout: 30 * time.Second,
	http:            just.Must(httph.New(httph.NewOptions())),
}
