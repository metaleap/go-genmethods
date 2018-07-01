package gentslices

import (
	"github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

var (
	// common var-names such as "i", "ok", "err", "this" etc
	ª = udevgogen.Vars

	// These "default `IGent`s" are a convenience offering in two ways:
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
}
