package gentenum

import (
	"fmt"

	gs "github.com/go-leap/dev/go/syn"
	"github.com/metaleap/go-gent"
)

// GentValidMethod generated a `Valid` method for enum type-defs,
// checking whether the value seems to be within the range of the
// known enumerants. It only supports enum type-defs whose enumerants
// are ordered in the source such that the smallest values appear first
// and the largest last, with all enumerant `const`s appearing together.
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

// GenerateTopLevelDecls implements github.com/metaleap/go-gent.IGent
func (this *GentValidMethod) GenerateTopLevelDecls(pkg *gent.Pkg, t *gent.Type) (tlDecls []gs.IEmit) {
	if t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {
		methodname, firstinvalid, firstname, lastname, firsthint, lasthint :=
			this.MethodName, this.IsFirstInvalid, t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1], "inclusive", "inclusive"
		if firstname == "_" {
			if len(t.Enumish.ConstNames) == 1 {
				return
			}
			firstinvalid, firstname = false, t.Enumish.ConstNames[1]
		}

		firstoperands, lastoperands := []gs.IEmit{gs.V.This, gs.N(firstname)}, []gs.IEmit{gs.V.This, gs.N(lastname)}
		firstoperator, lastoperator := gs.IEmit(gs.Geq(firstoperands...)), gs.IEmit(gs.Leq(lastoperands...))
		if firstinvalid {
			firstoperator, firsthint = gs.Gt(firstoperands...), "exclusive"
		}
		if this.IsLastInvalid {
			lastoperator, lasthint = gs.Lt(lastoperands...), "exclusive"
		}

		if methodname == "" {
			methodname = Defaults.Valid.MethodName
		}
		method := gs.Func(gs.V.This.T(gs.TrN("", t.Name)), methodname, gs.TrFunc(gs.TFunc(nil, gs.V.Ret.T(gs.TrpBool()))),
			gs.Set(gs.V.Ret, gs.And(firstoperator, lastoperator)),
		)
		method.DocCommentLines = append(method.DocCommentLines, fmt.Sprintf(this.DocComment, methodname, t.Name, firstname, firsthint, lastname, lasthint))
		tlDecls = append(tlDecls, method)
	}
	return
}
