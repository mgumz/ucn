package filter

// Limit returns the first limit runes if limit is larger than 0
func Limit(runes []rune, limit int) []rune {

	if limit > 0 && limit < len(runes) {
		runes = runes[:limit]
	}

	return runes
}
