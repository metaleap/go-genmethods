package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.Structs.Marshal.MethodName, Gents.Structs.Unmarshal.MethodName, Gents.Structs.Marshal.DocComment, Gents.Structs.Unmarshal.DocComment =
		DefaultMethodNameMarshal, DefaultMethodNameUnmarshal, DefaultDocCommentMarshal, DefaultDocCommentUnmarshal
}

type GentStructJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
	}
	Unmarshal struct {
		JsonMethodOpts
	}
}

func (me *GentStructJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Marshal.MethodName).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment).
		Code()
}

func (me *GentStructJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.MethodName, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment).
		Code()
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentStructJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.Expr.GenRef.Struct != nil {
		yield.Add(
			me.genMarshalMethod(ctx, t),
			me.genUnmarshalMethod(ctx, t),
		)
	}
	return
}
