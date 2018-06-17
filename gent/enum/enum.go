package gentenum

var (
	Defaults struct {
		Valid  GentValidMethod
		IsFoo  GentIsFooMethods
		String GentStringMethods
	}
)

func init() {
	Defaults.Valid.MethodName, Defaults.Valid.DocComment = "Valid", "%s returns whether the value of this `%s` is between `%s` (%s) and `%s` (%s)."
	Defaults.IsFoo.MethodNamePrefix, Defaults.IsFoo.DocComment = "Is", "%s returns whether the value of this `%s` equals `%s`."
	Defaults.String.Stringers = map[string]func(string) string{
		"String": nil,
	}
}
