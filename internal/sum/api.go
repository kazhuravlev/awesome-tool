package sum

import (
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
)

func Build(src source.Source, linksRules map[int][]source.Rule, linksFacts map[int]facts.Facts, linksChecks map[int]map[string][]rules.Error) Sum {
	return Sum{
		Version:     "1",
		GlobalRules: src.GlobalRulesEnabled,
		Groups:      src.Groups,
		Links:       src.Links,
		LinksRules:  linksRules,
		LinksFacts:  linksFacts,
		LinksChecks: linksChecks,
	}
}
