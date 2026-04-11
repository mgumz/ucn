package print

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"unicode"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

// isPUA reports whether r is in one of Unicode's three Private Use
// Areas. unicode.Co is the "Other, Private Use" general category and
// covers all three PUA blocks exactly; the stdlib unicode package is
// already linked into the binary by other call sites, so this costs
// nothing extra.
func isPUA(r rune) bool {
	return unicode.In(r, unicode.Co)
}

// ElasticTabstops prints given runes to io.Writer w
func ElasticTabstops(w io.Writer, entries []uc.Entry) {

	const (
		minWidth = 3
		tabWidth = 2
		padding  = 4
		padChar  = ' '
		flags    = 0
	)

	tw := tabwriter.NewWriter(w, minWidth, tabWidth, padding, padChar, flags)
	defer tw.Flush()

	// columns + format of each to print
	cols := strings.Join([]string{
		"%s",        // rune (pre-formatted below to let PUA glyphs pass through)
		"%04x",      // hex
		"%U",        //
		"%s",        // html-entity
		"\\\\u%04x", // json
		"%s\n"},     // name of rune
		"\t")

	for _, entry := range entries {
		r := entry.Rune()
		// %q escapes PUA runes as \uXXXX; bypass it for Private Use so an
		// installed Nerd Font (or any other PUA-using font) gets to render
		// the glyph. Non-PUA runes keep the %q quoting so control chars
		// like NUL don't corrupt the terminal.
		rfmt := "%q"
		if isPUA(r) {
			rfmt = "'%c'"
		}
		symbol := fmt.Sprintf(rfmt, r)
		fmt.Fprintf(tw, cols, symbol, r, r, runeToHTML(r), r, entry.Name)
	}
}
