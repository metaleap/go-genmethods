package gentenum

import (
	"fmt"
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentIsFooMethods generates a method `t.IsFoo() bool` for each
// enumerant `Foo` in enums, which equals-compares its receiver.
// (A hugely pointless code-gen, but its simplicity makes
// it a decent starter example for writing custom ones.)
type GentIsFooMethods struct {
	DocComment       string
	MethodNamePrefix string
	RenameEnumerant  func(string) string
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// If `t` is a suitable enum type-def, it returns a method `t.IsFoo() bool`
// for each enumerant `Foo` in `t`, which equals-compares its receiver.
func (this *GentIsFooMethods) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn) {
	if t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {
		trecv := TrNamed("", t.Name)
		tlDecls = make([]ISyn, 0, len(t.Enumish.ConstNames))
		for _, enumerant := range t.Enumish.ConstNames {
			if ren := enumerant; enumerant != "_" {
				if this.RenameEnumerant != nil {
					ren = this.RenameEnumerant(enumerant)
				}
				method := Func(V.This.Typed(trecv), this.MethodNamePrefix+ren, TdFunc(nil, V.Ret.Typed(T.Bool)),
					Set(V.Ret, Eq(V.This, N(enumerant))),
				)
				method.Doc.Add(fmt.Sprintf(this.DocComment, method.Name, t.Name, enumerant))
				tlDecls = append(tlDecls, method)
			}
		}
	}
	return
}
