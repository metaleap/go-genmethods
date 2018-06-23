package gentenum

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentIsFooMethods generates methods `YourEnumType.IsFoo() bool` for each enumerant `Foo`
// in enum type-defs, which equals-compares its receiver to the respective enumerant `Foo`.
// (A HIGHLY POINTLESS code-gen in real-world terms, except its exemplary simplicity makes
// it a handy starter-demo-sample-snippet-blueprint for writing new ones from scratch.)
//
// An instance with illustrative defaults is in `Defaults.IsFoo`.
type GentIsFooMethods struct {
	Disabled   bool
	DocComment gent.Str

	// eg `Is{e}` -> `IsMyOne`, `IsMyTwo`, etc.
	MethodName gent.Str

	// if set, renames the enumerant used for {e} in `MethodName`
	MethodNameRenameEnumerant func(string) string
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// If `t` is a suitable enum type-def, it returns a method `t.IsFoo() bool` for
// each enumerant `Foo` in `t`, which equals-compares its receiver to the enumerant.
func (this *GentIsFooMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if (!this.Disabled) && t.SeemsEnumish() {
		decls = make(Syns, 0, len(t.Enumish.ConstNames))
		for _, enumerant := range t.Enumish.ConstNames {
			if renamed := enumerant; enumerant != "_" {
				if this.MethodNameRenameEnumerant != nil {
					renamed = this.MethodNameRenameEnumerant(enumerant)
				}
				method := Fn(t.CodeGen.ThisVal, this.MethodName.With("{T}", t.Name, "{e}", renamed), &Sigs.NoneToBool,
					Set(V.Ret, Eq(V.This, N(enumerant))),
				)
				method.Doc.Add(this.DocComment.With(
					"{N}", method.Name,
					"{T}", t.Name,
					"{e}", enumerant,
				))
				decls.Add(method)
			}
		}
	}
	return
}
