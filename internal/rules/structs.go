package rules

import "github.com/kazhuravlev/awesome-tool/internal/sum"

type Check interface {
	Name() CheckName
	FactDeps() []sum.FactName
	Test(l sum.Link) (bool, []Error)
}

type Error string

type CheckName string
