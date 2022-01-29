package filter

import (
	"sort"

	fuzzy "github.com/lithammer/fuzzysearch/fuzzy"
	rn "golang.org/x/text/unicode/runenames"
)

// Fuzzy returns all runes matching the provided fuzzy filter
func Fuzzy(runes []rune, filter string) []rune {

	if filter == "" {
		return runes
	}

	names := []string{}
	for _, r := range runes {
		names = append(names, rn.Name(r))
	}

	matches := fuzzy.RankFindNormalizedFold(filter, names)
	sort.Sort(matches)

	mrunes := []rune{}
	for i := range matches {
		// Rank.OriginalIndex points to the index in "names",
		// which itself is the "named" representation of
		// "runes" … so the OriginalIndex should work in both
		// slices
		mrunes = append(mrunes, runes[matches[i].OriginalIndex])
	}

	return mrunes
}
