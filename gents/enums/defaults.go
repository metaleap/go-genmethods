package gentenums

import (
	"github.com/metaleap/go-gent"
)

var (
	// These `Defaults` are convenience offerings in two ways:
	// they illustrate usage of this package's individual `IGent`s' fields, and
	// they allow importers their own "defaults" base for less-noisy tweaking.
	// They are only initialized by this package, but not otherwise used by it.
	Defaults struct {
		IsFoo   GentIsFooMethods
		IsValid GentIsValidMethod
		List    GentListEnumerantsFunc
		String  GentStringMethods

		// contains pointers to all the above fields, in order
		All []gent.IGent
	}
)

func init() {
	Defaults.All = []gent.IGent{&Defaults.IsFoo, &Defaults.IsValid, &Defaults.List, &Defaults.String}

	Defaults.IsFoo.DocComment = "{N} returns whether the value of this `{T}` equals `{e}`."
	Defaults.IsFoo.MethodName = "Is{e}"

	Defaults.IsValid.DocComment = "{N} returns whether the value of this `{T}` is between `{fn}` ({fh}) and `{ln}` ({lh})."
	Defaults.IsValid.MethodName = "Valid"

	Defaults.List.DocComment = "{N} returns the `names` and `values` of all {n} well-known `{T}` enumerants."
	Defaults.List.FuncName = "Wellknown{T}{s}"

	defstr := &Defaults.String
	defstr.Stringers = []Stringer{
		{DocComment: "{N} implements the `fmt.Stringer` interface.", Name: "String",
			EnumerantRename: nil, ParseFuncName: "{T}From{str}", ParseAddErrlessVariantWithSuffix: "Or"},
		{DocComment: "{N} implements the `fmt.GoStringer` interface.", Name: "GoString",
			Disabled: true, ParseFuncName: "{T}From{str}", ParseAddErrlessVariantWithSuffix: "Or"},
	}
	defstr.DocComments.Parsers = "{N} returns the `{T}` represented by `{s}` (as returned by `{str}`, {caseSensitivity}), or an `error` if none exists."
	defstr.DocComments.ParsersErrlessVariant = "{N} is like `{p}` but returns `{fallback}` for bad inputs."
}
