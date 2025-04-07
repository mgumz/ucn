package main

import (
	"flag"
	"os"

	"github.com/mgumz/ucn/internal/pkg/print"
	uc "github.com/mgumz/ucn/internal/pkg/unicode"
)

type filterFunc func([]uc.Entry) []uc.Entry

func main() {

	printer := print.Func(print.ElasticTabstops)
	filters := []filterFunc{}

	initFlags(flag.CommandLine, &filters, &printer)

	flag.Parse()

	// input
	entries := uc.Entries()

	// transform
	for _, f := range filters {
		entries = f(entries)
	}

	// output
	printer(os.Stdout, entries)
}
