//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/ascii85"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	defaultUrl = "https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt"

	codeHeader = `// DO NOT EDIT: This file is generated!

package unicode

const ucEntries = ("" + `

	codeFooter = `  "");
`
	codeCommentFmt = "// %d bytes gzip+base64; the redundant, explicit parens are for https://golang.org/issue/18078"

	// https://golang.org/issue/18078 references commit
	//
	//   https://github.com/go-corelibs/x-text/commit/5c6cf4f9a2357d38515014cea8c488ed22bdab90#diff-7333055b57c2187f1e8cd500a210597595f9ba51bde4c8f58116536bc3fc50c6R205
	//
	// which introduces artificial parenthesis after 128kb. that matches roughly
	// the 62 lines for each block in the generated
	//
	//   https://cs.opensource.google/go/x/text/+/refs/tags/v0.24.0:unicode/runenames/tables15.0.0.go
	//
	// 64 just looks more … binary =) … and it should workaround the referred
	// issue good enough (tm).
	extraParensAfter = 64

	bytesPerLine = 72
)

func main() {

	url := flag.String("url", defaultUrl, "url to UnicodeData.txt")
	doFetchFromWeb := flag.Bool("fromWeb", false, "fetch UnicodeData.txt from web")
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
		r = req.Body
		defer req.Body.Close()
	}

	ibuf := bytes.NewBuffer(nil)

	// phase-1:
	// take the first two columns (delimiter is ";")
	// column-0: codepoint/symbol in hexadecimal
	// column-1: name / description
	// column-2+: rest
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			continue
		}

		name := parts[1]

		// Skip <control> or other category like entries. Column 10
		// has the "old name" - which is more useful than "<control>"
		// entries like. Also filtered: "<CJK Ideograph Extension A, First>"
		// which has no description/name at all.
		if (len(name) > 2) && (name[0] == '<') && (name[len(name)-1] == '>') {
			if len(parts) < 10 || parts[10] == "" {
				continue
			}
			name = parts[10]
		}

		rpl := strings.NewReplacer("(", "", ")", "")
		name = rpl.Replace(name)
		cp := parts[0]

		fmt.Fprintln(ibuf, cp, name)
	}

	// phase-2:
	// take the result of phase-1, gzip-9 it, base85 encode it
	obuf := bytes.NewBuffer(nil)
	a85w := ascii85.NewEncoder(obuf) // base85 is a bit more dense than base64
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
