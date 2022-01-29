package filter

import (
	"strings"

	rn "golang.org/x/text/unicode/runenames"
)

// Partial returns all runes matching the provided partial filter
func StartsWith(runes []rune, filter string) []rune {

	if filter == "" {
		return runes
	}

	filter = strings.ToLower(filter)
	matches := []rune{}
	for _, r := range runes {
		n := strings.ToLower(rn.Name(r))
		words := wsRE.Split(n, -1)
		for i := range words {
			if strings.HasPrefix(words[i], filter) {
				matches = append(matches, r)
				break
			}
		}
	}

	return matches
}
