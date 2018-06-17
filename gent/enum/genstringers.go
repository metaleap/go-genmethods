package gentenum

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentStringMethods struct {
	DocComment string
	Stringers  map[string]func(string) string
	Parsers    struct {
		OnePerStringer         bool
		OneUber                bool
		AddErrlessWithFallback bool
	}
}

func (this *GentStringMethods) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn) {
	if len(this.Stringers) > 0 && t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {
		tlDecls = make([]ISyn, 0, 2+len(t.Enumish.ConstNames)*len(this.Stringers)*3)
		for strname := range this.Stringers {
			if fns := this.genStringer(strname, t); fns != nil {
				tlDecls = append(tlDecls, fns)
			}
			if this.Parsers.OnePerStringer {
				if fnp := this.genParser(strname, t); fnp != nil {
					tlDecls = append(tlDecls, fnp)
				}
			}
		}
		if this.Parsers.OneUber {
		}
	}
	return
}

func (this *GentStringMethods) genStringer(strName string, t *gent.Type) (method *SynFunc) {
	method = Fn(t.CodeGen.MethodRecvVal, strName, &Sigs.NoneToString,
		Set(V.Ret, L("ret")),
	)
	return
}

func (this *GentStringMethods) genParser(strName string, t *gent.Type) (method *SynFunc) {
	return
}
