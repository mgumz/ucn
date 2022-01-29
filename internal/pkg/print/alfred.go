package print

import (
	"encoding/json"
	"fmt"
	"io"

	rn "golang.org/x/text/unicode/runenames"
)

// https://www.alfredapp.com/help/workflows/inputs/script-filter/json/
//
// {"items": [
//     {
//         "uid": "desktop",
//         "type": "file",
//         "title": "Desktop",
//         "subtitle": "~/Desktop",
//         "arg": "~/Desktop",
//         "autocomplete": "Desktop",
//         "icon": {
//             "type": "fileicon",
//             "path": "~/Desktop"
//         }
//     }
// ]}

type asf struct {
	Items []asfi `json:"items"`
}

type asfi struct {
	UID       string    `json:"uid,omitempty"`
	Type      string    `json:"type,omitempty"`
	Title     string    `json:"title"`
	SubTitle  string    `json:"subtitle"`
	Arg       string    `json:"arg,omitempty"`
	Variables jsonEntry `json:"variables,omitempty"`
	Text      *struct {
		Copy      string `json:"copy,omitempty"`
		LargeType string `json:"largetype,omitempty"`
	} `json:"text,omitempty"`
}

// AlfredJSON prints given runes to io.Writer in
// the JSON flavor suitable for alfredapp.com
// (where it is useful in Alfred Workflows)
func AlfredJSON(w io.Writer, runes []rune) {
	result := asf{}
	for _, r := range runes {
		n := rn.Name(r)
		h := runeToHTML(r)
		u := fmt.Sprintf(baseWWWURL, r)
		item := asfi{
			// Why no UID? "ucn" yields the proper order and
			// without UID Alfred is respecting the order as
			// it was defined, see:
			// https://www.alfredapp.com/help/workflows/inputs/script-filter/json/#uid
			//UID:      fmt.Sprintf("%U", r),

			// Arg: send to "output" of Alfred's workflow
			Arg: string(r),

			// Variables:
			// https://www.deanishe.net/post/2018/10/workflow/environment-variables-in-alfred/
			Variables: jsonEntry{
				Symbol: string(r),
				CP:     fmt.Sprintf("%U", r),
				HTML:   h,
				JSON:   fmt.Sprintf("\\u%x", r),
				WWW:    u,
			},

			//
			Title:    fmt.Sprintf("%c - %s", r, n),
			SubTitle: fmt.Sprintf("%U - %s - %s", r, h, u),
		}
		// item.Text.Copy = item.Title
		// item.Text.LargeType = item.Title
		result.Items = append(result.Items, item)
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("  ", "  ")
	_ = enc.Encode(&result)
}
