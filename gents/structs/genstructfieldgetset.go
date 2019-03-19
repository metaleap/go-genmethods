package gentstructs

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultDocCommentGet = ""
	DefaultMethodNameGet = "StructFieldsGet"
	DefaultDocCommentSet = ""
	DefaultMethodNameSet = "StructFieldsSet"
)

func init() {
	Gents.StructFieldsGetSet.DocCommentGet, Gents.StructFieldsGetSet.MethodNameGet = DefaultDocCommentGet, DefaultMethodNameGet
	Gents.StructFieldsGetSet.DocCommentSet, Gents.StructFieldsGetSet.MethodNameSet = DefaultDocCommentSet, DefaultMethodNameSet
}

type GentStructFieldsGetSet struct {
	gent.Opts

	DontGen struct {
		Getters bool
		Setters bool
	}
	DocCommentGet gent.Str
	MethodNameGet string
	DocCommentSet gent.Str
	MethodNameSet string
}

func (me *GentStructFieldsGetSet) genGetMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.MethodNameGet).
		Arg("name", T.String).
		Rets(ˇ.R.OfType(T.Empty.Interface), ˇ.Ok.OfType(T.Bool)).
		Doc().
		Code()
}

func (me *GentStructFieldsGetSet) genSetMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.MethodNameSet).
		Args(T.String.N("name"), ˇ.V.OfType(T.Empty.Interface)).
		Rets(T.Bool.N("okName"), T.Bool.N("okType")).
		Doc().
		Code()
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentStructFieldsGetSet) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if tstruct := t.Expr.GenRef.Struct; tstruct != nil {
		if !me.DontGen.Getters {
			yield.Add(me.genGetMethod(ctx, t))
		}
		if !me.DontGen.Setters {
			yield.Add(me.genSetMethod(ctx, t))
		}
	}
	return
}
