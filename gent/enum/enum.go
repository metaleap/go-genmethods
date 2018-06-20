package gentenum

var (
	Defaults struct {
		Valid  GentValidMethod
		IsFoo  GentIsFooMethods
		String GentStringMethods
	}
)

func init() {
	Defaults.Valid.MethodName, Defaults.Valid.DocComment = "Valid", "{N} returns whether the value of this `{T}` is between `{fn}` ({fh}) and `{ln}` ({lh})."
	Defaults.IsFoo.MethodNamePrefix, Defaults.IsFoo.DocComment = "Is", "{N} returns whether the value of this `{T}` equals `{e}`."
	Defaults.String.Parsers.FuncName, Defaults.String.Parsers.OnePerStringer, Defaults.String.Stringers = "{T}From{s}", true, map[string]func(string) string{
		"String": nil,
	}
}
