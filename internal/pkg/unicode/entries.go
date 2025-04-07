package unicode

import (
	"bufio"
	"compress/gzip"
	"encoding/ascii85"
	"strconv"
	"strings"
)

type Entry struct {
	Symbol string
	Name   string
}

func (e *Entry) Rune() rune {
	n, _ := strconv.ParseUint(e.Symbol, 16, 32)
	return rune(n)
}

func Entries() []Entry {

	sr := strings.NewReader(ucEntries)
	dec85 := ascii85.NewDecoder(sr)
	gzr, _ := gzip.NewReader(dec85)

	entries := []Entry{}

	scanner := bufio.NewScanner(gzr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		entries = append(entries, Entry{Symbol: parts[0], Name: parts[1]})
	}

	return entries
}
