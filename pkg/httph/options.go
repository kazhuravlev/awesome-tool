package httph

import (
	"net/http"
)

//go:generate options-gen -from-struct=Options -defaults-from=var
type Options struct {
	client *http.Client

	maxEquivRedirects int
}

var defaultOptions = Options{
	client:            http.DefaultClient,
	maxEquivRedirects: 3,
}
