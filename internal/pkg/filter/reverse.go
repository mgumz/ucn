package filter

import "slices"

// Reverse reverses the order of the given runes
func Reverse(runes []rune, reverse bool) []rune {
	if reverse {
		slices.Reverse(runes)
	}
	return runes
}
