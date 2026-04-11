//go:build ignore

// Generator for the Nerd Fonts glyph metadata blob.
//
// Upstream file: glyphnames.json from https://github.com/ryanoasis/nerd-fonts
// License of the metadata: MIT (see upstream LICENSE).
//
// Version pinning: nerdFontsRef is a git tag on the upstream repo. Bumping
// it is an intentional act — change the constant, re-run `go generate`,
// review the diff, and note the bump in the commit message.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/ascii85"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

const (
	// Bump intentionally. Releases: https://github.com/ryanoasis/nerd-fonts/releases
	nerdFontsRef = "v3.4.0"

	defaultUrl = "https://raw.githubusercontent.com/ryanoasis/nerd-fonts/" +
		nerdFontsRef + "/glyphnames.json"

	codeHeader = `// DO NOT EDIT: This file is generated!

package nerdfont

const nfEntries = ("" + `

	codeFooter = `  "");
`
	codeCommentFmt = "// %d bytes gzip+base64; the redundant, explicit parens are for https://golang.org/issue/18078"

	// See the long comment in internal/pkg/unicode/data_gen.go for the
	// rationale behind the periodic extra parentheses.
	extraParensAfter = 64

	bytesPerLine = 72
)

type nfGlyph struct {
	Char string `json:"char"`
	Code string `json:"code"`
}

func main() {

	url := flag.String("url", defaultUrl, "url to glyphnames.json")
	doFetchFromWeb := flag.Bool("fromWeb", false, "fetch glyphnames.json from web")
	oFile := flag.String("out", "", "file to write to")
	flag.Parse()

	r := io.Reader(os.Stdin)
	owriter := os.Stdout
	if *oFile != "" {
		var err error
		if owriter, err = os.Create(*oFile); err != nil {
			fmt.Fprintf(os.Stderr, "error creating file %q: %s\n", *oFile, err)
			os.Exit(13)
		}
	}

	if *doFetchFromWeb {
		req, err := http.Get(*url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error fetching %q: %s\n", *url, err)
			os.Exit(1)
		}
		if req.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "error fetching %q: HTTP %s\n", *url, req.Status)
			os.Exit(1)
		}
		r = req.Body
		defer req.Body.Close()
	}

	// phase-1:
	// decode glyphnames.json leniently via json.RawMessage: the top-level
	// map mixes glyph entries ({char, code}) with a METADATA entry whose
	// shape is different. Decode each value on its own; silently skip any
	// that don't fit the glyph shape. Emit one line per glyph:
	// "<code> <UPPER-NAME-WITH-SPACES>".
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(r).Decode(&raw); err != nil {
		fmt.Fprintf(os.Stderr, "error decoding json: %s\n", err)
		os.Exit(2)
	}

	glyphs := make(map[string]nfGlyph, len(raw))
	for name, rm := range raw {
		if name == "METADATA" {
			continue
		}
		var g nfGlyph
		if err := json.Unmarshal(rm, &g); err != nil {
			continue
		}
		if g.Code == "" {
			continue
		}
		glyphs[name] = g
	}

	names := make([]string, 0, len(glyphs))
	for name := range glyphs {
		names = append(names, name)
	}
	sort.Strings(names) // deterministic output

	rpl := strings.NewReplacer("-", " ", "_", " ", "(", "", ")", "")

	ibuf := bytes.NewBuffer(nil)
	for _, name := range names {
		g := glyphs[name]
		display := strings.ToUpper(rpl.Replace(name))
		cp := strings.ToUpper(g.Code)
		fmt.Fprintln(ibuf, cp, display)
	}

	// phase-2:
	// take the result of phase-1, gzip-9 it, base85 encode it
	obuf := bytes.NewBuffer(nil)
	a85w := ascii85.NewEncoder(obuf)
	gzw, _ := gzip.NewWriterLevel(a85w, gzip.BestCompression)
	io.Copy(gzw, ibuf)
	gzw.Close()
	a85w.Close()

	sb := strings.Builder{}
	sb.Grow(80)

	// phase-3: construct data.go
	comment := fmt.Sprintf(codeCommentFmt, obuf.Len())
	fmt.Fprintln(owriter, codeHeader, comment)

	for i := 1; obuf.Len() > 0; i++ {
		sb.Reset()
		sb.WriteString("  \"")
		line := obuf.Next(bytesPerLine)
		for _, b := range line {
			switch b {
			case '"':
				sb.WriteString("\\\"")
			case '\\':
				sb.WriteString("\\\\")
			default:
				sb.WriteByte(b)
			}
		}
		if (i % extraParensAfter) == 0 {
			sb.WriteString("\") + (")
		} else {
			sb.WriteString("\" +")
		}

		fmt.Fprintln(owriter, sb.String())
	}
	fmt.Fprint(owriter, codeFooter)
}
