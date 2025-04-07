package filter

import (
	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// Limit returns the first limit runes if limit is larger than 0
func Limit(entries []uc.Entry, limit int) []uc.Entry {

	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	return entries
}
