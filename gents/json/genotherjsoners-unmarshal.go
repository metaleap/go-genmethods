package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func (me *GentTypeJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var code ISyn
	var selftype *TypeRef
	if t.Expr.GenRef.Struct != nil {
		selftype = t.G.Tª
		code = me.genUnmarshalStruct(ctx, func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef })
	} else {
		selftype = t.G.T
	}
	return selftype.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				B.Panic.Of(t.Name),
			), Else(
				code,
			)),
		)
}

func (me *GentTypeJsonMethods) genUnmarshalStruct(ctx *gent.Ctx, field func() (ISyn, *TypeRef)) (code Syns) {
	return
}
