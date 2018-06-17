package gentenum

import (
	"fmt"

	. "github.com/go-leap/dev/go/syn"
	"github.com/metaleap/go-gent"
)

// GentValidMethod generates a `Valid` method for enum type-defs that
// checks whether the receiver value seems to be within the range of the
// known enumerants. It's only correct for enum type-defs whose enumerants
// are ordered in the source such that the smallest values appear first,
// the largest last, and with all enumerant `const`s appearing together.
type GentValidMethod struct {
	// defaults to Defaults.Valid.DocComment
	DocComment string

	// defaults to Defaults.Valid.MethodName
	MethodName string

	// if `true`, generate gt instead of geq
	IsFirstInvalid bool

	// if `true`, generate lt instead of leq
	IsLastInvalid bool
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// It returns at most one method if `t` is a suitable enum type-def.
func (this *GentValidMethod) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn) {
	if t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {
		methodname, firstinvalid, firstname, lastname, firsthint, lasthint :=
			this.MethodName, this.IsFirstInvalid, t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1], "inclusive", "inclusive"
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

		if methodname == "" {
			methodname = Defaults.Valid.MethodName
		}
		method := Func(V.This.Typed(TrNamed("", t.Name)), methodname, TdFunc(nil, V.Ret.Typed(T.Bool)),
			Set(V.Ret, And(firstoperator, lastoperator)),
		)
		method.DocCommentLines = append(method.DocCommentLines, fmt.Sprintf(this.DocComment, methodname, t.Name, firstname, firsthint, lastname, lasthint))
		tlDecls = append(tlDecls, method)
	}
	return
}
