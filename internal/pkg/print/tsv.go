package print

import (
	"encoding/csv"
	"fmt"
	"io"

	rn "golang.org/x/text/unicode/runenames"
)

// TSV prints given runes to io.Writer w
func TSV(w io.Writer, runes []rune) {

	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	defer cw.Flush()

	for _, r := range runes {
		cw.Write([]string{
			string(r),
			fmt.Sprintf("%04x", r),
			fmt.Sprintf("%U", r),
			runeToHTML(r),
			fmt.Sprintf("\\\\u%04x", r),
			rn.Name(r),
		})
	}
}
