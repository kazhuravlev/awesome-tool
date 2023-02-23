package sum

import "github.com/kazhuravlev/just"

func getLinksCountRecursive(groups ...Group) int {
	return just.Sum(just.SliceMap(groups, func(g Group) int {
		return len(g.Links)
	})...)
}
