package sum

import (
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
)

type Sum struct {
	// Version is version of .sum file
	Version string
	Groups  []Group
}

type Group struct {
	SrcGroup source.Group
	Groups   []Group
	Links    []Link
	// Contains count of links for this group and all children groups
	LinksCountRecursive int
	IsPresentInResult bool
}

type Link struct {
	SrcLink source.Link
	Rules   []Rule
	Facts   facts.Facts
}

type Rule struct {
	SrcRule source.Rule
	Checks  []Check
}

type Check struct {
	Name         rules.CheckName
	IsTestPassed bool
	Errors       []rules.Error
}
