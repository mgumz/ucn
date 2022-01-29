package main

import (
	"flag"
	"strconv"

	"github.com/mgumz/ucn/internal/pkg/filter"
	"github.com/mgumz/ucn/internal/pkg/print"
)

const (
	flagVersion = "print version and exit"

	flagOfmtEtab   = "use Elastic Tabstops (default)"
	flagOfmtTSV    = "use tab-separated values"
	flagOfmtJSON   = "use JSON to print the symbols"
	flagOfmtAlfred = "use special JSON format ready to be used in http://alfredapp.com"

	flagFilterPF         = "filter all symbol names for a partial match on <value>"
	flagFilterFF         = "fuzzy filter all symbol names for <value>"
	flagFilterStartsWith = "starts-with filter all symbol names for <value>"
	flagFilterReverse    = "reverse list of symbols"
	flagFilterLimit      = "limit results to <value> symbols"
)

func initFlags(fs *flag.FlagSet, filters *[]filterFunc, printer *print.Func) {

	fs.BoolFunc("v", flagVersion, printVersion)

	// -ofmt.xyz
	fs.BoolFunc("ofmt.etab", flagOfmtEtab, func(_ string) error { *printer = print.ElasticTabstops; return nil })
	fs.BoolFunc("ofmt.tsv", flagOfmtTSV, func(_ string) error { *printer = print.TSV; return nil })
	fs.BoolFunc("ofmt.alfred", flagOfmtAlfred, func(_ string) error { *printer = print.AlfredJSON; return nil })
	fs.BoolFunc("ofmt.json", flagOfmtJSON, func(_ string) error { *printer = print.JSON; return nil })

	// filters
	fs.Func("filter.partial", flagFilterPF, func(f string) error {
		*filters = append(*filters, func(runes []rune) []rune {
			return filter.Partial(runes, f)
		})
		return nil
	})

	fs.Func("filter.fuzzy", flagFilterFF, func(f string) error {
		*filters = append(*filters, func(runes []rune) []rune {
			return filter.Fuzzy(runes, f)
		})
		return nil
	})

	fs.Func("filter.starts-with", flagFilterStartsWith, func(f string) error {
		*filters = append(*filters, func(runes []rune) []rune {
			return filter.StartsWith(runes, f)
		})
		return nil
	})

	fs.Func("reverse", flagFilterReverse, func(f string) error {
		*filters = append(*filters, func(runes []rune) []rune {
			return filter.Reverse(runes, true)
		})
		return nil
	})

	fs.Func("limit", flagFilterLimit, func(l string) error {
		limit, err := strconv.ParseInt(l, 10, 32)
		if err != nil {
			return err
		}
		*filters = append(*filters, func(runes []rune) []rune {
			return filter.Limit(runes, int(limit))
		})
		return nil
	})
}
