package print

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// ElasticTabstops prints given runes to io.Writer w
func ElasticTabstops(w io.Writer, entries []uc.Entry) {

	tw := tabwriter.NewWriter(w, 3, 2, 4, ' ', 0)
	defer tw.Flush()

	// columns + format of each to print
	cols := strings.Join([]string{
		"%q",        // rune
		"%04x",      // hex
		"%U",        //
		"%s",        // html-entity
		"\\\\u%04x", // json
		"%s\n"},     // name of rune
		"\t")

	for _, entry := range entries {
		r := entry.Rune()
		fmt.Fprintf(tw, cols, r, r, r, runeToHTML(r), r, entry.Name)
	}
}
