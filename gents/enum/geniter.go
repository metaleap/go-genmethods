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
	}
	return
}

func (this *GentIterateFuncs) genIterWithCallback(t *gent.Type, withNameArg bool, withValArg bool) *SynFunc {
	trcallback := TrFunc(TdFunc(nil))
	if withNameArg {
		trcallback.Func.Args.Add("", T.String)
	}
	if withValArg {
		trcallback.Func.Args.Add("", TrNamed("", t.Enumish.BaseType))
	}
	return Fn(NoMethodRecv, this.IterWithCallbackFuncName.With("{T}", t.Name), TdFunc(NTs("onEnumerant", trcallback)),
		K.Ret,
	)
}
