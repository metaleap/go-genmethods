// Package gentenums provides `gent.IGent` code-gens of `func`s related to "enum-ish
// type-defs". Most of them expect and assume enum type-defs whose enumerants are
// ordered in the source such that the numerically smallest value appears first,
// the largest one last, with all enumerant `const`s appearing next to each other.
package gentenums

import (
	"github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

var (
	// common var-names such as "i", "ok", "err", "this" etc.
	ˇ = udevgogen.Vars

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IsFoo     GentIsFooMethods
		IsValid   GentIsValidMethod
		List      GentListEnumerantsFunc
		Stringers GentStringersMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)

func init() {
	Gents.All = gent.Gents{&Gents.IsFoo, &Gents.IsValid, &Gents.List, &Gents.Stringers}
}
