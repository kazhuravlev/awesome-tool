package sum

import (
	"errors"
	"fmt"

	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
)

// gExtractors contains all available etractors. This is like a registry of all
// extractors.
var gExtractors = make(map[FactName]FactExtractor)

// gExtractorsDeps contains fact dependencies.
var gExtractorsDeps = make(map[FactName][]FactName)

// RegisterExtractor will add extractor in global registry.
func RegisterExtractor(extractor FactExtractor) error {
	if just.MapContainsKey(gExtractors, extractor.Name()) {
		return errors.New("fact extractor already exists")
	}

	for _, dep := range extractor.Deps() {
		if !just.MapContainsKey(gExtractors, dep) {
			return fmt.Errorf("extractor '%s' not found", dep)
		}
	}

	gExtractors[extractor.Name()] = extractor
	gExtractorsDeps[extractor.Name()] = extractor.Deps()

	return nil
}

func MustRegisterExtractor(extractor FactExtractor) {
	if err := RegisterExtractor(extractor); err != nil {
		panic("register extractor: " + err.Error())
	}
}

func GatherFacts(link source.Link) (*Link, error) {
	resLink := Link{
		SrcLink: link,
		Facts:   LinkFacts{},
	}

	for i := range gExtractors {
		gExtractors[i].Extract(&resLink)
	}

	return &resLink, nil
}

type FactExtractor interface {
	// Name of extractor
	Name() FactName
	// Deps of this extractor
	Deps() []FactName
	// Implementation of extractor
	Extract(*Link)
}
