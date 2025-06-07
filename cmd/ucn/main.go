package main

import (
	"flag"
	"os"

	"github.com/mgumz/ucn/internal/pkg/filter"
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

	// treat regular args as list of codepoints to "filter" on. default: all.
	// side-effect: the user can "ask" on what the name of of a code point is
	// by just providing the hex-code. (saves use of grep)
	if len(flag.Args()) == 0 {
		goto transform
	}

	// thus, we _prepend_ the filter for the code points to the defined
	// other filters.
	filters = append([]filterFunc{
		func(entries []uc.Entry) []uc.Entry {
			return filter.CodeInSlice(entries, flag.Args())
		},
	}, filters...)

transform:

	for _, f := range filters {
		entries = f(entries)
	}

	// output
	printer(os.Stdout, entries)
}
