package filter

import (
	"sort"

	fuzzy "github.com/lithammer/fuzzysearch/fuzzy"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// Fuzzy returns all runes matching the provided fuzzy filter
func Fuzzy(entries []uc.Entry, filter string) []uc.Entry {

	if filter == "" {
		return entries
	}

	names := []string{}
	for i := range entries {
		names = append(names, entries[i].Name)
	}

	matches := fuzzy.RankFindNormalizedFold(filter, names)
	sort.Sort(matches)

	filtered := []uc.Entry{}
	for i := range matches {
		// Rank.OriginalIndex points to the index in "names",
		// which itself is the "named" representation of
		// "runes" … so the OriginalIndex should work in both
		// slices
		filtered = append(filtered, entries[matches[i].OriginalIndex])
	}

	return filtered
}
