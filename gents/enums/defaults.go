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

	Gents.IsFoo.DocComment = "{N} returns whether the value of this `{T}` equals `{e}`."
	Gents.IsFoo.MethodName = "Is{e}"

	Gents.IsValid.DocComment = "{N} returns whether the value of this `{T}` is between `{fn}` ({fh}) and `{ln}` ({lh})."
	Gents.IsValid.MethodName = "Valid"

	Gents.List.DocComment = "{N} returns the `names` and `values` of all {n} well-known `{T}` enumerants."
	Gents.List.FuncName = "Wellknown{T}{s}"

	defstr := &Gents.String
	defstr.Stringers = []StringMethodOpts{
		{DocComment: "{N} implements the `fmt.Stringer` interface.", Name: "String",
			EnumerantRename: nil, ParseFuncName: "{T}From{str}", ParseErrless: gent.Variant{Add: false, NameOrSuffix: "Or"}},
		{DocComment: "{N} implements the `fmt.GoStringer` interface.", Name: "GoString",
			Disabled: true, ParseFuncName: "{T}From{str}", ParseErrless: gent.Variant{Add: false, NameOrSuffix: "Or"}},
	}
	defstr.DocComments.Parsers = "{N} returns the `{T}` represented by `{s}` (as returned by `{str}`, {caseSensitivity}), or an `error` if none exists."
	defstr.DocComments.ParsersErrlessVariant = "{N} is like `{p}` but returns `{fallback}` for bad inputs."
}
