package gentenum

var (
	Defaults struct {
		Valid  GentValidMethod
		IsFoo  GentIsFooMethods
		String GentStringMethods
		Iters  GentIterateFuncs
	}
)

func init() {
	Defaults.Valid.MethodName, Defaults.Valid.DocComment = "Valid", "{N} returns whether the value of this `{T}` is between `{fn}` ({fh}) and `{ln}` ({lh})."
	Defaults.IsFoo.MethodNamePrefix, Defaults.IsFoo.DocComment = "Is", "{N} returns whether the value of this `{T}` equals `{e}`."
	Defaults.String.Stringers = []Stringer{
		{Name: "String", EnumerantRename: nil, ParseFuncName: "{T}From{s}", ParseAddErrlessVariantWithSuffix: "Or"},
	}
	Defaults.Iters.IterWithCallbackFuncName, Defaults.Iters.EnumerantsFuncName = "ForEachWellknown{T}", "Wellknown{T}s"
}
