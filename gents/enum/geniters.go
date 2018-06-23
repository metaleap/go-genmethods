package gentenum

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentIterateFuncs struct {
	EnumerantsFuncName            gent.Str
	IterWithCallbackFuncName      gent.Str
	NoEnumerantNameArgInCallback  bool
	NoEnumerantValueArgInCallback bool
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentIterateFuncs) GenerateTopLevelDecls(t *gent.Type) (tlDecls Syns) {
	if t.SeemsEnumish() {
		if withnamearg, withvalarg := !this.NoEnumerantNameArgInCallback, !this.NoEnumerantValueArgInCallback; this.IterWithCallbackFuncName != "" && (withnamearg || withvalarg) {
			tlDecls.Add(this.genIterWithCallback(t, withnamearg, withvalarg))
		}
		if this.EnumerantsFuncName != "" {
			tlDecls.Add(this.genIterEnumerants(t))
		}
	}
	return
}

func (this *GentIterateFuncs) genIterEnumerants(t *gent.Type) *SynFunc {
	names, values := make(Syns, 0, len(t.Enumish.ConstNames)), make(Syns, 0, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if enumerant != "_" {
			names, values = append(names, L(enumerant)), append(values, NT(enumerant, t.CodeGen.Ref))
		}
	}

	fn := Fn(NoMethodRecv, this.EnumerantsFuncName.With("{T}", t.Name), TdFunc(nil, NT("names", T.Sl.Strings), NT("values", TrSlice(t.CodeGen.Ref))),
		Set(C(N("names"), N("values")), C(L(names), L(values))),
	)
	return fn
}

func (this *GentIterateFuncs) genIterWithCallback(t *gent.Type, withNameArg bool, withValArg bool) *SynFunc {
	trcallback := TrFunc(TdFunc(nil))
	if withNameArg {
		trcallback.Func.Args.Add("", T.String)
	}
	if withValArg {
		trcallback.Func.Args.Add("", t.CodeGen.Ref)
	}
	n, fn := N("onEnumerant"), Fn(NoMethodRecv, this.IterWithCallbackFuncName.With("{T}", t.Name), TdFunc(NTs("onEnumerant", trcallback)))
	for _, enumerant := range t.Enumish.ConstNames {
		if enumerant != "_" {
			fn.Add(Call(n, L(enumerant), N(enumerant)))
		}
	}
	return fn
}
