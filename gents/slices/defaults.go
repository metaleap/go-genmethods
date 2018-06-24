package gentslices

import (
	"github.com/metaleap/go-gent"
)

var (
	// These `Defaults` are convenience offerings in two ways:
	// they illustrate usage of this package's individual `IGent`s' fields, and
	// they allow importers their own "defaults" base for less-noisy tweaking.
	// They are only initialized by this package, but not otherwise used by it.
	Defaults struct {
		IndexOf GentIndexMethods

		// contains pointers to all the above fields, in order
		All []gent.IGent
	}
)

func init() {
	Defaults.All = []gent.IGent{&Defaults.IndexOf}

	defidx := &Defaults.IndexOf
	defidx.IndexOf.Name, defidx.IndicesOf.Name, defidx.IndexLast.Name =
		"Index", "Indices", "LastIndex"
	defidx.IndexOf.PredicateVariation.Name, defidx.IndicesOf.PredicateVariation.Name, defidx.IndexLast.PredicateVariation.Name =
		"IndexFunc", "IndicesFunc", "LastIndexFunc"
}
