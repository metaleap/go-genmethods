package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultIsFooDocComment = "{N} returns whether the value of this `{T}` equals `{e}`."
	DefaultIsFooMethodName = "Is{e}"
)

func init() {
	Gents.IsFoo.DocComment, Gents.IsFoo.MethodName = DefaultIsFooDocComment, DefaultIsFooMethodName
}

// GentIsFooMethods generates methods `YourEnumType.IsFoo() bool` for each enumerant `Foo`
// in enum type-defs, which equals-compares its receiver to the respective enumerant `Foo`.
// (A HIGHLY POINTLESS code-gen in real-world terms, except its exemplary simplicity makes
// it a handy starter-demo-sample-snippet-blueprint for writing new ones from scratch.)
//
// An instance with illustrative defaults is in `Gents.IsFoo`.
type GentIsFooMethods struct {
	gent.Opts

	DocComment gent.Str
	// eg `Is{e}` -> `IsMyOne`, `IsMyTwo`, etc.
	MethodName gent.Str

	// if set, renames the enumerant used for {e} in `MethodName`
	MethodNameRenameEnumerant func(string) string
}

func (this *GentIsFooMethods) genIsFooMethod(t *gent.Type, methodName string, enumerant string) *SynFunc {
	return t.Gen.ThisVal.Method(methodName).Sig(&Sigs.NoneToBool).
		Doc(this.DocComment.With("{N}", methodName, "{T}", t.Name, "{e}", enumerant)).
		Code(
			V.R.SetTo(V.This.Eq(N(enumerant))),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// If `t` is a suitable enum type-def, it returns a method `t.IsFoo() bool` for
// each enumerant `Foo` in `t`, which equals-compares its receiver to the enumerant.
func (this *GentIsFooMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsEnumish() {
		decls = make(Syns, 0, len(t.Enumish.ConstNames))
		for _, enumerant := range t.Enumish.ConstNames {
			if renamed := enumerant; enumerant != "_" {
				if this.MethodNameRenameEnumerant != nil {
					renamed = this.MethodNameRenameEnumerant(enumerant)
				}
				decls.Add(this.genIsFooMethod(t, this.MethodName.With("{T}", t.Name, "{e}", renamed), enumerant))
			}
		}
	}
	return
}
