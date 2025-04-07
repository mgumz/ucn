package print

import (
	"encoding/csv"
	"fmt"
	"io"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// TSV prints given runes to io.Writer w
func TSV(w io.Writer, entries []uc.Entry) {

	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	defer cw.Flush()

	for _, entry := range entries {
		r := entry.Rune()
		cw.Write([]string{
			string(r),
			fmt.Sprintf("%04x", r),
			fmt.Sprintf("%U", r),
			runeToHTML(r),
			fmt.Sprintf("\\\\u%04x", r),
			entry.Name,
		})
	}
}
