package rules

import (
	"fmt"

	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/source"
)

type CheckResponseStatusEq struct {
	ExpectedStatus int
}

func (c CheckResponseStatusEq) Name() CheckName {
	return "response:status-eq"
}

func (c CheckResponseStatusEq) FactDeps() []facts.FactName {
	return []facts.FactName{facts.FactResponse}
}

func (c CheckResponseStatusEq) Test(link source.Link, linkFacts facts.Facts) CheckResult {
	if linkFacts.Data.Response.StatusCode != c.ExpectedStatus {
		return CheckResult{
			CheckName: c.Name(),
			IsPassed:  false,
			Errors: []Error{
				Error(fmt.Sprintf(
					"response status code is '%d', but should be '%d'",
					linkFacts.Data.Response.StatusCode,
					c.ExpectedStatus,
				)),
			},
		}
	}

	return CheckResult{
		CheckName: c.Name(),
		IsPassed:  true,
		Errors:    nil,
	}
}

type CheckGithubStarsMin struct {
	MinStars int
}

func (c CheckGithubStarsMin) Name() CheckName {
	return "github-repo:stars-min"
}

func (c CheckGithubStarsMin) FactDeps() []facts.FactName {
	return []facts.FactName{facts.FactGithub}
}

func (c CheckGithubStarsMin) Test(link source.Link, linkFacts facts.Facts) CheckResult {
	if linkFacts.Data.Github.StargazersCount < c.MinStars {
		return CheckResult{
			CheckName: c.Name(),
			IsPassed:  false,
			Errors: []Error{
				Error(fmt.Sprintf(
					"repository has not enough stars '%d'. Required at least '%d' stars",
					linkFacts.Data.Github.StargazersCount,
					c.MinStars,
				)),
			},
		}
	}

	return CheckResult{
		CheckName: c.Name(),
		IsPassed:  true,
		Errors:    nil,
	}
}
