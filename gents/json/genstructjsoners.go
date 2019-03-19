package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.Structs.DocCommentMarshal, Gents.Structs.DocCommentUnmarshal =
		DefaultDocCommentMarshal, DefaultDocCommentUnmarshal
}

type GentStructJsonMethods struct {
	gent.Opts

	DocCommentMarshal   string
	DocCommentUnmarshal string
}

func (this *GentStructJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.T.Method("MarshalJSON").Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(this.DocCommentMarshal).
		Code()
}

func (this *GentStructJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method("UnmarshalJSON", ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(this.DocCommentUnmarshal).
		Code()
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStructJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.Expr.GenRef.Struct != nil {
		yield.Add(
			this.genMarshalMethod(ctx, t),
			this.genUnmarshalMethod(ctx, t),
		)
	}
	return
}
