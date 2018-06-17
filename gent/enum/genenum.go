package gentenum

import (
	gs "github.com/go-leap/dev/go/syn"
	"github.com/metaleap/go-gent"
)

var (
	Defaults struct {
		Valid  GentValidMethod
		IsFoo  GentIsFooMethods
		String GentStringMethods
	}
)

func init() {
	Defaults.Valid.MethodName = "Valid"
}

type GentStringMethods struct {
}

// GentValidMethod works for enumish `type`s whose
// enumerants are ordered such that the smallest
// values appear first and the largest last.
type GentValidMethod struct {
	// Defaults.Valid.MethodName
	MethodName string

	// first of the enumerants
	IsFirstInvalid bool

	// last of the enumerants
	IsLastInvalid bool
}

func (this *GentValidMethod) GenerateTopLevelDecls(pkg *gent.Pkg, t *gent.Type) (tlDecls []gs.IEmit) {
	if methodname := this.MethodName; t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {
		isfirstinvalid, namefirst, namelast := this.IsFirstInvalid, t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1]
		if namefirst == "_" {
			if len(t.Enumish.ConstNames) == 1 {
				return
			}
			isfirstinvalid, namefirst = false, t.Enumish.ConstNames[1]
		}

		opsg, opsl := []gs.IEmit{gs.V.This, gs.N(namefirst)}, []gs.IEmit{gs.V.This, gs.N(namelast)}
		opg, opl := gs.IEmit(gs.Geq(opsg...)), gs.IEmit(gs.Leq(opsl...))
		if this.IsLastInvalid {
			opg = gs.Gt(opsg...)
		}
		if isfirstinvalid {
			opl = gs.Lt(opsl...)
		}

		if methodname == "" {
			methodname = Defaults.Valid.MethodName
		}
		method := gs.Func(gs.V.This.T(gs.TrN("", t.Name)), methodname, gs.TrFunc(gs.TFunc(nil, gs.V.Ret.T(gs.TrpBool()))),
			gs.Set(gs.V.Ret, gs.And(opg, opl)),
		)
		tlDecls = append(tlDecls, method)
	}
	return
}

type GentIsFooMethods struct {
}
