package gentenums

import (
	"strconv"
	"strings"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultListDocComment = "{N} returns the `names` and `values` of all {n} well-known `{T}` enumerants."
	DefaultListFuncName   = "Wellknown{T}{s}"
)

func init() {
	Gents.List.DocComment, Gents.List.FuncName = DefaultListDocComment, DefaultListFuncName
}

// GentListEnumerantsFunc generates a
// `func WellknownFoos() ([]string, []Foo)`
// for each enum type-def `Foo`.
//
// An instance with illustrative defaults is in `Gents.List`.
type GentListEnumerantsFunc struct {
	gent.Opts

	DocComment gent.Str
	// eg. "Wellknown{T}{s}" with `{T}` for type name and
	// `{s}` for pluralization suffix (either "s" or "es")
	FuncName gent.Str
}

func (this *GentListEnumerantsFunc) genListEnumerantsFunc(t *gent.Type, funcName string, enumerantNames Syns, enumerantValues Syns) *SynFunc {
	return Func(funcName).Ret("names", T.Sl.Strings).Ret("values", TrSlice(t.G.T)).
		Doc(this.DocComment.With("N", funcName, "T", t.Name, "n", strconv.Itoa(len(enumerantNames)))).
		Code(
			Names("names", "values").SetTo(Lits(enumerantNames, enumerantValues)),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentListEnumerantsFunc) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsEnumish() {
		names, values := make(Syns, 0, len(t.Enumish.ConstNames)), make(Syns, 0, len(t.Enumish.ConstNames))
		for _, enumerant := range t.Enumish.ConstNames {
			if enumerant != "_" {
				names, values = append(names, L(enumerant)), append(values, NT(enumerant, t.G.T))
			}
		}
		pluralsuffix := "s"
		if strings.HasSuffix(t.Name, "s") {
			pluralsuffix = "es"
		}
		fname := this.FuncName.With("T", t.Name, "s", pluralsuffix)
		if strings.HasSuffix(fname, "ys") {
			fname = fname[:len(fname)-2] + "ies"
		}
		decls.Add(this.genListEnumerantsFunc(t, fname, names, values))
	}
	return
}
