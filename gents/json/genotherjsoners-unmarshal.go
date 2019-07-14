package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func (me *GentTypeJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				B.Panic.Of(t.Name),
			), Else()),
		)
}
