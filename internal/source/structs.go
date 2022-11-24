package source

import "github.com/kazhuravlev/just"

type RuleName string
type GroupName string
type CheckName string

type Source struct {
	Version            string     `yaml:"version"`
	Rules              []Rule     `yaml:"rules"`
	GlobalRulesEnabled []RuleName `yaml:"global-rules-enabled"`
	Groups             []Group    `yaml:"groups"`
	Links              []Link     `yaml:"links"`
}

type Check struct {
	Check CheckName `yaml:"check"`
	Args  []string  `yaml:"args"`
}

type Rule struct {
	Name   RuleName `yaml:"name"`
	Title  string   `yaml:"title"`
	Checks []Check  `yaml:"checks"`
}

type Group struct {
	Name         GroupName               `yaml:"name"`
	Title        string                  `yaml:"title"`
	Description  just.NullVal[string]    `yaml:"description"`
	Group        just.NullVal[GroupName] `yaml:"group,omitempty"`
	RulesEnabled []RuleName              `yaml:"rules-enabled"`
	RulesIgnored []RuleName              `yaml:"rules-ignored"`
	AlwaysShown  bool                    `yaml:"always-shown"`
}

type Link struct {
	URL          string      `yaml:"url"`
	Title        string      `yaml:"title"`
	RulesEnabled []RuleName  `yaml:"rules-enabled"`
	RulesIgnored []RuleName  `yaml:"rules-ignored"`
	Groups       []GroupName `yaml:"groups"`
}
