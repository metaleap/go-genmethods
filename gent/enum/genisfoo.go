package gentenum

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentIsFooMethods generates methods `YourEnumType.IsFoo() bool` for each enumerant `Foo`
// in enum type-defs, which equals-compares its receiver to the respective enumerant `Foo`.
// (A highly pointless code-gen in real-world terms, except its exemplary simplicity
// makes it a handy starter demo sample snippet for writing new ones from scratch.)
type GentIsFooMethods struct {
	DocComment       gent.Str
	MethodNamePrefix string
	RenameEnumerant  func(string) string
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// If `t` is a suitable enum type-def, it returns a method `t.IsFoo() bool` for
// each enumerant `Foo` in `t`, which equals-compares its receiver to the enumerant.
func (this *GentIsFooMethods) GenerateTopLevelDecls(t *gent.Type) (tlDecls Syns) {
	if t.SeemsEnumish() {
		tlDecls = make(Syns, 0, len(t.Enumish.ConstNames))
		for _, enumerant := range t.Enumish.ConstNames {
			if ren := enumerant; enumerant != "_" {
				if this.RenameEnumerant != nil {
					ren = this.RenameEnumerant(enumerant)
				}
				method := Fn(t.CodeGen.ThisVal, this.MethodNamePrefix+ren, &Sigs.NoneToBool,
					Set(V.Ret, Eq(V.This, N(enumerant))),
				)
				method.Doc.Add(this.DocComment.With(
					"{N}", method.Name,
					"{T}", t.Name,
					"{e}", enumerant,
				))
				tlDecls.Add(method)
			}
		}
	}
	return
}
