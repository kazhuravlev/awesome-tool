package source

import (
	"errors"
	"fmt"

	"github.com/kazhuravlev/awesome-tool/internal/errorsh"
	"github.com/kazhuravlev/just"
)

func ParseFile(filename string) (*Source, error) {
	obj, err := unmarshalFilename(filename)
	if err != nil {
		return nil, errorsh.Wrap(err, "unmarshal filename")
	}

	if obj.Version != "1" {
		return nil, errors.New("unknown source file version")
	}

	return obj, nil
}

func Validate(obj Source) error {
	rulesMap := make(map[RuleName]struct{}, len(obj.Rules))
	for i := range obj.Rules {
		rule := &obj.Rules[i]

		if just.MapContainsKey(rulesMap, rule.Name) {
			return fmt.Errorf("rule '%s' is duplicated", rule.Name)
		}

		if len(rule.Checks) == 0 {
			return fmt.Errorf("rule '%s' has no checks", rule.Name)
		}

		// FIXME: check the check names and args

		rulesMap[rule.Name] = struct{}{}
	}

	for _, ruleName := range obj.GlobalRulesEnabled {
		if !just.MapContainsKey(rulesMap, ruleName) {
			return errorsh.Newf("unknown global-enabled-rule '%s'", ruleName)
		}
	}

	groupsMap := make(map[GroupName]struct{}, len(obj.Groups))
	for i := range obj.Groups {
		group := &obj.Groups[i]

		if just.MapContainsKey(groupsMap, group.Name) {
			return fmt.Errorf("group '%s' is duplicated", group.Name)
		}

		if parentGroupName, ok := group.Group.ValueOk(); ok {
			if parentGroupName == group.Name {
				return fmt.Errorf("group '%s' refers to itself", group.Name)
			}

			if !just.MapContainsKey(groupsMap, parentGroupName) {
				return fmt.Errorf("group '%s' has unknown parent '%s'", group.Name, parentGroupName)
			}
		}

		for _, ruleName := range just.SliceChain(group.RulesEnabled, group.RulesIgnored) {
			if !just.MapContainsKey(rulesMap, ruleName) {
				return fmt.Errorf("group '%s' refers to not exists rule '%s'", group.Name, ruleName)
			}
		}

		groupsMap[group.Name] = struct{}{}
	}

	linksMap := make(map[string]struct{}, len(obj.Links))
	for i := range obj.Links {
		link := &obj.Links[i]

		if just.MapContainsKey(linksMap, link.URL) {
			return fmt.Errorf("link '%s' is duplicated", link.URL)
		}

		for i := range link.Groups {
			if !just.MapContainsKey(groupsMap, link.Groups[i]) {
				return fmt.Errorf("link '%s' refers to not exists group '%s'", link.URL, link.Groups[i])
			}
		}

		for _, ruleName := range just.SliceChain(link.RulesEnabled, link.RulesIgnored) {
			if !just.MapContainsKey(rulesMap, ruleName) {
				return fmt.Errorf("link '%s' refers to not exists rule '%s'", link.URL, ruleName)
			}
		}

		linksMap[link.URL] = struct{}{}
	}

	return nil
}
