package sum

import (
	"github.com/kazhuravlev/awesome-tool/internal/facts"
	"github.com/kazhuravlev/awesome-tool/internal/rules"
	"github.com/kazhuravlev/awesome-tool/internal/source"
	"github.com/kazhuravlev/just"
)

const Version = "1"

func Build(
	src source.Source,
	linksRules map[int][]source.Rule,
	linksFacts map[int]facts.Facts,
	linksChecks map[int]map[string]rules.CheckResult,
) Sum {
	group2links := make(map[source.GroupName][]source.Link, len(src.Groups))
	for _, link := range src.Links {
		for _, groupName := range link.Groups {
			group2links[groupName] = append(group2links[groupName], link)
		}
	}

	groups := make([]Group, 0, len(src.Groups))
	for _, group := range src.Groups {
		// We need only root groups, so we skip groups, that has parent group.
		if group.Group.Valid {
			continue
		}

		groups = append(groups, handleGroup(src, group, linksRules, linksFacts, linksChecks))
	}

	return Sum{
		Version: Version,
		Groups:  groups,
	}
}

func handleGroup(
	src source.Source,
	g source.Group,
	linksRules map[int][]source.Rule,
	linksFacts map[int]facts.Facts,
	linksChecks map[int]map[string]rules.CheckResult,
) Group {
	groupLinks := make([]Link, 0, len(src.Links))
	for linkIdx, link := range src.Links {
		if !just.SliceContainsElem(link.Groups, g.Name) {
			continue
		}

		linkRules := make([]Rule, len(linksRules[linkIdx]))
		for i, rule := range linksRules[linkIdx] {
			checks := make([]Check, len(rule.Checks))
			for i, checkRawString := range rule.Checks {
				checkresults := linksChecks[linkIdx][checkRawString]
				// TODO: this structs duplicate rules.CheckResults.
				checks[i] = Check{
					Name:         checkresults.CheckName,
					IsTestPassed: checkresults.IsPassed,
					Errors:       checkresults.Errors,
				}
			}

			linkRules[i] = Rule{
				SrcRule: rule,
				Checks:  checks,
			}
		}

		groupLinks = append(groupLinks, Link{
			SrcLink: link,
			Rules:   linkRules,
			Facts:   linksFacts[linkIdx],
		})
	}

	var childGroups []Group
	for _, group := range src.Groups {
		if group.Group.Val != g.Name {
			continue
		}

		childGroups = append(childGroups, handleGroup(src, group, linksRules, linksFacts, linksChecks))
	}

	linksCountRecursive := len(groupLinks) + getLinksCountRecursive(childGroups...)
	return Group{
		SrcGroup:            g,
		Groups:              childGroups,
		Links:               groupLinks,
		LinksCountRecursive: linksCountRecursive,
		IsPresentInResult:   linksCountRecursive > 0 || g.AlwaysShown,
	}
}
