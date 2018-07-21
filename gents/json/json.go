package gentjson

import (
	"github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

var (
	// common var-names such as "i", "ok", "err", "this" etc
	ˇ = udevgogen.Vars

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		Enums GentEnumJsonMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)

func init() {
	Gents.All = gent.Gents{&Gents.Enums}
}
