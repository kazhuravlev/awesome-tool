package sum

import (
	"net/url"
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
