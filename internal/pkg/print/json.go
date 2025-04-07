package print

import (
	"encoding/json"
	"fmt"
	"io"

	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

type jsonEntry struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name,omitempty"`
	CP     string `json:"cp"`
	HTML   string `json:"html"`
	JSON   string `json:"json"`
	WWW    string `json:"www"`
}

func JSON(w io.Writer, entries []uc.Entry) {

	j := []jsonEntry{}
	for _, entry := range entries {
		r := entry.Rune()
		j = append(j, jsonEntry{
			Symbol: string(r),
			Name:   entry.Name,
			CP:     fmt.Sprintf("%U", r),
			HTML:   runeToHTML(r),
			JSON:   fmt.Sprintf("\\u%04x", r),
			WWW:    fmt.Sprintf(baseWWWURL, r),
		})
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("  ", " ")
	enc.Encode(j)
}
