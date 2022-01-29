package print

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	rn "golang.org/x/text/unicode/runenames"
)

// ElasticTabstops prints given runes to io.Writer w
func ElasticTabstops(w io.Writer, runes []rune) {

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

	for _, r := range runes {
		fmt.Fprintf(tw, cols, r, r, r, runeToHTML(r), r, rn.Name(r))
	}
}
