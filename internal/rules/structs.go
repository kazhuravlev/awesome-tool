package rules

import (
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
)

type Check interface {
	Name() CheckName
	FactDeps() []facts.FactName
	Test(l source.Link, lFacts facts.Facts) CheckResult
}

type Error string

type CheckName string

type CheckResult struct {
	CheckName CheckName
	IsPassed  bool
	Errors    []Error
}

type RulesResults struct {
	Reports map[source.RuleName]map[CheckName]bool
}

func (r RulesResults) ReportByRule() map[source.RuleName]bool {
	res := make(map[source.RuleName]bool, len(r.Reports))
	for ruleName, check2status := range r.Reports {
		res[ruleName] = !just.SliceContainsElem(just.MapGetValues(check2status), false)
	}

	return res
}

// IsOK reutrns true when all rules is in completed state.
func (r RulesResults) IsOK() bool {
	// NOTE: no rules - means that all is ok already.
	if len(r.Reports) == 0 {
		return true
	}

	return !just.SliceContainsElem(just.MapGetValues(r.ReportByRule()), false)
}
