package app

import (
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/kazhuravlev/awesome-tool/pkg/httph"
	"github.com/kazhuravlev/just"
)

type Encoder interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error

	MarshalFile(string, any) error
	UnmarshalFile(string, any) error
}

//go:generate options-gen -from-struct=Options -defaults-from=var

type Options struct {
	responseTimeout time.Duration

	githubClient *github.Client

	http       *httph.Client
	maxWorkers int
	sumEncoder Encoder
}

var defaultOptions = Options{
	responseTimeout: 30 * time.Second,
	http:            just.Must(httph.New(httph.NewOptions())),
	maxWorkers:      1,
	sumEncoder:      YamlEncoder{},
}
