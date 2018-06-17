package gentenum

import (
	"fmt"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentValidMethod generates a `Valid` method for enum type-defs, which
// checks whether the receiver value seems to be within the range of the
// known enumerants. It is only correct for enum type-defs whose enumerants
// are ordered in the source such that the numerically smallest values appear
// first, the largest ones last, with all enumerant `const`s appearing together.
type GentValidMethod struct {
	DocComment     string
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// It returns at most one method if `t` is a suitable enum type-def.
func (this *GentValidMethod) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn) {
	if t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {
		firstinvalid, firstname, lastname, firsthint, lasthint :=
			this.IsFirstInvalid, t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1], "inclusive", "inclusive"
		if firstname == "_" {
			if len(t.Enumish.ConstNames) == 1 {
				return
			}
			firstinvalid, firstname = false, t.Enumish.ConstNames[1]
		}

		var firstoperator, lastoperator ISyn = Geq(V.This, N(firstname)), Leq(V.This, N(lastname))
		if firstinvalid {
			firstoperator, firsthint = Gt(V.This, N(firstname)), "exclusive"
		}
		if this.IsLastInvalid {
			lastoperator, lasthint = Lt(V.This, N(lastname)), "exclusive"
		}

		method := Fn(t.CodeGen.MethodRecvVal, this.MethodName, &Sigs.NoneToBool,
			Set(V.Ret, And(firstoperator, lastoperator)),
		)
		method.Doc.Add(fmt.Sprintf(this.DocComment, method.Name, t.Name, firstname, firsthint, lastname, lasthint))
		tlDecls = append(tlDecls, method)
	}
	return
}
