package gentenum

import (
	gs "github.com/go-leap/dev/go/syn"
	"github.com/metaleap/go-gent"
)

var (
	GentValidMethodDefaultMethodName = "Valid"
)

type GentStringMethods struct {
}

type GentValidMethod struct {
	// defaults to GentValidMethodDefaultMethodName
	MethodName string
}

func (this *GentValidMethod) GenerateTopLevelDecls(pkg *gent.Pkg, t *gent.Type) (tlDecls []gs.IEmit) {
	if methodname := this.MethodName; t.Enumish.Potentially && len(t.Enumish.ConstantNames) > 0 {
		if methodname == "" {
			methodname = GentValidMethodDefaultMethodName
		}
		method := gs.Func(gs.V.This.T(gs.TrN("", t.Name)), methodname, gs.TrFunc(gs.TFunc(nil, gs.V.Ret.T(gs.TrpBool()))))
		tlDecls = append(tlDecls, method)
	}
	return
}

type GentIsEnumerantMethods struct {
}
