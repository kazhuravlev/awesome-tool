package rules

import "github.com/kazhuravlev/awesome-tool/internal/sum"

type Check interface {
	Name() CheckName
	Test(sum.Link) (bool, []Error)
}

type Error struct{}

type CheckName string
