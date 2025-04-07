package filter

import (
	"slices"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// Reverse reverses the order of the given runes
func Reverse(entries []uc.Entry, reverse bool) []uc.Entry {
	if reverse {
		slices.Reverse(entries)
	}
	return entries
}
