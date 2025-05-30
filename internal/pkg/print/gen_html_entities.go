//go:build ignore

// this code generates html_entities.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// {
//  "&AElig": { "codepoints": [198], "characters": "\u00C6" },
//  "&AElig;": { "codepoints": [198], "characters": "\u00C6" },

type EntityList map[string]struct {
	CPs   []int  `json:"codepoints"`
	Chars string `json:"characters"`
}

func main() {

	url := "https://html.spec.whatwg.org/entities.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error fetching", url, ":", err)
		os.Exit(1)
		return
	}

	defer resp.Body.Close()

	entities := EntityList{}
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&entities); err != nil {
		fmt.Fprintln(os.Stderr, "error decoding", url, ":", err)
		os.Exit(2)
		return
	}

	rmap := map[rune]string{}

	for e, v := range entities {
		// we have "&AMP" and "&AMP;", lets only use "&AMP;"
		if !strings.HasSuffix(e, ";") {
			continue
		}

		r, s := utf8.DecodeRuneInString(v.Chars)
		// we also skip all entries with multiples runes
		if s != len(v.Chars) {
			continue
		}

		// we want to prefer lower-case entities over Upper-Case ones
		// (&amp; trumps &AMP;)
		me, exists := rmap[r]
		if exists {
			r2, _ := utf8.DecodeRuneInString(me[1:])
			if unicode.IsUpper(r2) {
				continue
			}
		}

		rmap[r] = e
	}

	// prepare stable output order
	rorder := make([]rune, 0, len(rmap))
	for r := range rmap {
		rorder = append(rorder, r)
	}
	slices.Sort(rorder)

	fmt.Printf(`//go:generate go run gen_html_entities.go > html_entities.go

// code generated by go generate; DO NOT EDIT.

package print

// based upon %q
var rune2entity = map[rune]string {
`, url)
	for _, r := range rorder {
		fmt.Printf("	%d: %q,\n", r, rmap[r])
	}
	fmt.Println("}")
}
