package print

import (
	"strconv"
)

func runeToHTML(r rune) string {

	if e, exists := rune2entity[r]; exists {
		return e
	}

	return "&#" + strconv.FormatInt(int64(r), 16) + ";"
}
