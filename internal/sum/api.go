package sum

import (
	"net/url"

	"github.com/kazhuravlev/awesome-tool/internal/source"
)

func GatherFacts(link source.Link) (*Link, error) {
	resLink := Link{
		SrcLink: link,
		Facts:   LinkFacts{},
	}

	extractors := []factExtractor{}
	for i := range extractors {
		extractors[i](&resLink)
	}

	return &resLink, nil
}

type factExtractor func(*Link)

func extractUrl(link *Link) {
	u, err := url.Parse(link.SrcLink.URL)
	if err != nil {
		return
	}

	link.Facts.Url.Data = u
}

func extractGithub(link *Link) {
	if !link.Facts.Url.HasData {
		return
	}

	u := link.Facts.Url.Data

	link.Facts.Url.Data = u
}
