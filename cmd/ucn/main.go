package main

import (
	"flag"
	"os"

	rt "golang.org/x/text/unicode/rangetable"
	rn "golang.org/x/text/unicode/runenames"

	"github.com/mgumz/ucn/internal/pkg/print"
)

type filterFunc func([]rune) []rune

func main() {

	printer := print.Func(print.ElasticTabstops)
	filters := []filterFunc{}

	initFlags(flag.CommandLine, &filters, &printer)

	flag.Parse()

	// input
	runes := initRunes()

	// transform
	for _, f := range filters {
		runes = f(runes)
	}

	// output
	printer(os.Stdout, runes)
}

func initRunes() []rune {

	runes := []rune{}
	table := rt.Assigned(rn.UnicodeVersion)
	rt.Visit(table, func(r rune) {
		runes = append(runes, r)
	})
	return runes
}
