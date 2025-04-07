package filter

import (
	"strings"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// Partial returns all runes matching the provided partial filter
func StartsWith(entries []uc.Entry, filter string) []uc.Entry {

	if filter == "" {
		return entries
	}

	filter = strings.ToLower(filter)
	matches := []uc.Entry{}
	for _, e := range entries {
		n := strings.ToLower(e.Name)
		words := wsRE.Split(n, -1)
		for i := range words {
			if strings.HasPrefix(words[i], filter) {
				matches = append(matches, e)
				break
			}
		}
	}

	return matches
}
