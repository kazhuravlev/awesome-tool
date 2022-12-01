package sum

import (
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
)

type Sum struct {
	// Version is version of .sum file
	Version string
	// GlobalRules - all rules, that will applied to all links. This field does
	// not contains all rules, that registered in source file. Only those which
	// will be applied globally.
	GlobalRules []source.RuleName
	Groups      []source.Group
	Links       []source.Link
	// LinksRUles contains a map of link idx => list of rules that will be
	// applied concrete for this link.
	LinksRules map[int][]source.Rule
	// LinksChecks contains map of source.check => list of check errors for
	// each link idx.
	LinksChecks map[int]map[string][]rules.Error
	// LinksFacts is a map of link idx => facts
	LinksFacts map[int]facts.Facts
}
