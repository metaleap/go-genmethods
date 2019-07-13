package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
	"github.com/metaleap/go-gent/gents/enums"
)

func init() {
	Gents.Enums.Marshal.MethodName, Gents.Enums.Unmarshal.MethodName, Gents.Enums.Marshal.DocComment, Gents.Enums.Unmarshal.DocComment, Gents.Enums.StringerToUse =
		DefaultMethodNameMarshal, DefaultMethodNameUnmarshal, DefaultDocCommentMarshal, DefaultDocCommentUnmarshal, &gentenums.Gents.Stringers.All[0]
}

type GentEnumJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
	}
	Unmarshal struct {
		JsonMethodOpts
	}
	StringerToUse *gentenums.StringMethodOpts
}

func (me *GentEnumJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.T.Method(me.Marshal.MethodName).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment).
		Code(
			ˇ.R.Set(T.SliceOf.Bytes.From(ctx.Import("strconv").C("Quote", Self.C(me.StringerToUse.Name)))),
		)
}

func (me *GentEnumJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.MethodName, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment).
		Code(
			Var(ˇ.S.Name, T.String, nil),
			Tup(ˇ.S, ˇ.Err).Set(ctx.Import("strconv").C("Unquote", T.String.From(ˇ.V))),
			If(ˇ.Err.Eq(B.Nil), Then(
				Var(ˇ.T.Name, t.G.T, nil),
				Tup(ˇ.T, ˇ.Err).Set(N(me.StringerToUse.Parser.FuncName.With("T", t.Name, "str", me.StringerToUse.Name)).Of(ˇ.S)),
				If(ˇ.Err.Eq(B.Nil), Then(
					Self.Deref().Set(ˇ.T),
				)),
			)),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentEnumJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsEnumish() {
		yield.Add(
			me.genMarshalMethod(ctx, t),
			me.genUnmarshalMethod(ctx, t),
		)
	}
	return
}
