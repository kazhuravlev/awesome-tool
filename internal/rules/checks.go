package rules

import (
	"fmt"
	"github.com/kazhuravlev/awesome-tool/internal/sum"
)

type CheckResponseStatusEq struct {
	ExpectedStatus int
}

func (c CheckResponseStatusEq) Name() CheckName {
	return "response:status-eq"
}

func (c CheckResponseStatusEq) FactDeps() []sum.FactName {
	return []sum.FactName{sum.FactResponse}
}

func (c CheckResponseStatusEq) Test(link sum.Link) (bool, []Error) {
	if link.Facts.Response.StatusCode != c.ExpectedStatus {
		return false, []Error{
			Error(fmt.Sprintf(
				"response status code is '%d', but should be '%d'",
				link.Facts.Response.StatusCode,
				c.ExpectedStatus,
			)),
		}
	}
	return true, nil
}

type CheckGithubStarsMin struct {
	MinStars int
}

func (c CheckGithubStarsMin) Name() CheckName {
	return "github-repo:stars-min"
}

func (c CheckGithubStarsMin) FactDeps() []sum.FactName {
	return []sum.FactName{sum.FactGithub}
}

func (c CheckGithubStarsMin) Test(link sum.Link) (bool, []Error) {
	if link.Facts.Github.StargazersCount < c.MinStars {
		return false, []Error{
			Error(fmt.Sprintf(
				"repository has not enough stars '%d'. Required at least '%d' stars",
				link.Facts.Github.StargazersCount,
				c.MinStars,
			)),
		}
	}
	return true, nil
}
