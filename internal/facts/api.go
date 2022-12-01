package facts

import (
	"context"
	"errors"
	"fmt"
	"github.com/kazhuravlev/awesome-tool/internal/errorsh"

	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
)

// gExtractors contains all available extractors. This is like a registry of all
// extractors.
var gExtractors = make(map[FactName]FactExtractor)

// gExtractorsOrdering contains order of fact extractors. This is executing
// order.
var gExtractorsOrdering []FactName

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
	gExtractorsOrdering = append(gExtractorsOrdering, extractor.Name())

	return nil
}

func MustRegisterExtractor(extractor FactExtractor) {
	if err := RegisterExtractor(extractor); err != nil {
		panic("register extractor: " + err.Error())
	}
}

func GatherFacts(ctx context.Context, link source.Link) (*Facts, error) {
	facts, err := gatherFactsLink(ctx, link)
	if err != nil {
		return nil, errorsh.Wrap(err, "gather facts about link")
	}

	return facts, nil
}

func gatherFactsLink(ctx context.Context, link source.Link) (*Facts, error) {
	var facts Data
	readyMap := make(map[FactName]bool, len(gExtractors))

ExtractCycle:
	for _, factName := range gExtractorsOrdering {
		extractor := gExtractors[factName]

		for _, dep := range extractor.Deps() {
			if !readyMap[dep] {
				continue ExtractCycle
			}
		}

		ok, err := extractor.Extract(ctx, link, &facts)
		if err != nil {
			return nil, errorsh.Wrap(err, "extract fact")
		}

		readyMap[extractor.Name()] = ok
	}

	return &Facts{
		Collected: readyMap,
		Data:      facts,
	}, nil
}
