package print

import (
	"encoding/json"
	"fmt"
	"io"

	rn "golang.org/x/text/unicode/runenames"
)

type jsonEntry struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name,omitempty"`
	CP     string `json:"cp"`
	HTML   string `json:"html"`
	JSON   string `json:"json"`
	WWW    string `json:"www"`
}

func JSON(w io.Writer, runes []rune) {

	entries := []jsonEntry{}
	for _, r := range runes {
		entries = append(entries, jsonEntry{
			Symbol: string(r),
			Name:   rn.Name(r),
			CP:     fmt.Sprintf("%U", r),
			HTML:   runeToHTML(r),
			JSON:   fmt.Sprintf("\\u%04x", r),
			WWW:    fmt.Sprintf(baseWWWURL, r),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("  ", " ")
	enc.Encode(entries)
}
