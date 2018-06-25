package gentslices

import (
	"github.com/metaleap/go-gent"
)

var (
	// These "default `IGent`s" are convenience offerings in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IndexOf GentIndexMethods

		// contains pointers to all the above fields, in order
		All []gent.IGent
	}
)

func init() {
	Gents.All = []gent.IGent{&Gents.IndexOf}

	defidx := &Gents.IndexOf
	defidx.IndexOf.Name, defidx.IndicesOf.Name, defidx.IndexLast.Name, defidx.Contains.Name, defidx.IndicesOf.Disabled, defidx.IndexLast.Disabled, defidx.Contains.Disabled =
		"Index", "Indices", "LastIndex", "Contains", true, true, false
	defidx.IndexOf.Predicate.NameOrSuffix, defidx.IndicesOf.Predicate.NameOrSuffix, defidx.IndexLast.Predicate.NameOrSuffix =
		"Func", "Func", "Func"
}
