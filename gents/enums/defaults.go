package gentenums

import (
	"github.com/metaleap/go-gent"
)

var (
	// These "default `IGent`s" are convenience offerings in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IsFoo   GentIsFooMethods
		IsValid GentIsValidMethod
		List    GentListEnumerantsFunc
		String  GentStringMethods

		// contains pointers to all the above fields, in order
		All []gent.IGent
	}
)

func init() {
	Gents.All = []gent.IGent{&Gents.IsFoo, &Gents.IsValid, &Gents.List, &Gents.String}
}
