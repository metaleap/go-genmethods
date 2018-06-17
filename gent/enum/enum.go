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
}
