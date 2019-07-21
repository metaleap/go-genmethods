package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
	"github.com/metaleap/go-gent/gents/enums"
)

func init() {
	Gents.EnumishTypes.Marshal.Name, Gents.EnumishTypes.Unmarshal.Name, Gents.EnumishTypes.Marshal.DocComment, Gents.EnumishTypes.Unmarshal.DocComment, Gents.EnumishTypes.StringerToUse =
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

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentEnumJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsEnumish() {
		if gennormal, genpanic := me.Marshal.genWhat(t); gennormal || genpanic {
			yield.Add(me.genMarshalMethod(ctx, t, genpanic))
		}
		if gennormal, genpanic := me.Unmarshal.genWhat(t); gennormal || genpanic {
			yield.Add(me.genUnmarshalMethod(ctx, t, genpanic))
		}
	}
	return
}

func (me *GentEnumJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	return t.G.T.Method(me.Marshal.Name).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment.With("N", me.Marshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				B.Panic.Of(t.Name),
			), Else(
				ˇ.R.Set(T.SliceOf.Bytes.From(ctx.Import("strconv").C("Quote", Self.C(me.StringerToUse.Name)))),
			)),
		)
}

func (me *GentEnumJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				Then(B.Panic.Of(t.Name)),
			), Else(
				Var(ˇ.S.Name, T.String, nil),
				Tup(ˇ.S, ˇ.Err).Set(ctx.Import("strconv").C("Unquote", T.String.From(ˇ.V))),
				If(ˇ.Err.Eq(B.Nil), Then(
					Var(ˇ.T.Name, t.G.T, nil),
					Tup(ˇ.T, ˇ.Err).Set(N(me.StringerToUse.Parser.FuncName.With("T", t.Name, "str", me.StringerToUse.Name)).Of(ˇ.S)),
					If(ˇ.Err.Eq(B.Nil), Then(
						Self.Deref().Set(ˇ.T),
					)),
				)),
			)),
		)
}
