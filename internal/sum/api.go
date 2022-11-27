package sum

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

func gatherFactsLink(ctx context.Context, link source.Link) (*LinkFacts, error) {
	var facts LinkFacts

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

	return &facts, nil
}

func GatherFacts(ctx context.Context, sourceObj source.Source) (*Sum, error) {
	rulesMap := make(map[source.RuleName]source.Rule)
	for _, rule := range sourceObj.Rules {
		rulesMap[rule.Name] = rule
	}

	// get global enabled rules
	globalRules := just.SliceMap(sourceObj.GlobalRulesEnabled, func(rn source.RuleName) source.Rule {
		return rulesMap[rn]
	})

	links := make([]Link, len(sourceObj.Links))
	for i, link := range sourceObj.Links {
		// FIXME: implement group-level rules
		// FIXME: implement link-level rules
		facts, err := gatherFactsLink(ctx, link)
		if err != nil {
			return nil, err
		}

		links[i] = Link{
			SrcLink: link,
			Facts:   *facts,
			Rules:   globalRules,
		}
	}

	return &Sum{
		Version: sourceObj.Version,
		Rules:   sourceObj.Rules,
		Groups:  sourceObj.Groups,
		Links:   links,
	}, nil
}
