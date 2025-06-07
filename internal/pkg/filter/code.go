package filter

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

func CodeInSlice(entries []uc.Entry, codes []string) []uc.Entry {

	if len(codes) == 0 {
		return entries
	}

	ncodes := make([]string, 0, len(codes))
	r := strings.NewReplacer(
		"U+", "",
		"&#", "",
		"\\U", "",
		"0X", "",
	)
	for _, code := range codes {
		code = strings.TrimSpace(code)
		if utf8.RuneCountInString(code) == 1 { // '¡'
			r, _ := utf8.DecodeRuneInString(code)
			if r != utf8.RuneError {
				code = fmt.Sprintf("%.4x", r)
			}
		}

		// &#00a1 or 0x00a1 or whatever

		code = strings.ToUpper(code) // entry.Symbol is upper-case hex
		code = r.Replace(code)

		// ensure uniqeness after normalization
		if slices.Contains(ncodes, code) {
			continue
		}

		ncodes = append(ncodes, code)
	}

	matches := []uc.Entry{}
	for _, entry := range entries {
		for _, code := range ncodes {
			if strings.Contains(entry.Symbol, code) {
				matches = append(matches, entry)
			}
		}
	}

	return matches
}
