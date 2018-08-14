package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
	"github.com/metaleap/go-gent/gents/enums"
)

const (
	DefaultDocCommentMarshal   = "MarshalJSON implements the Go standard library's `encoding/json.Marshaler` interface."
	DefaultDocCommentUnmarshal = "UnmarshalJSON implements the Go standard library's `encoding/json.Unmarshaler` interface."
)

func init() {
	Gents.Enums.DocCommentMarshal, Gents.Enums.DocCommentUnmarshal, Gents.Enums.StringerToUse =
		DefaultDocCommentMarshal, DefaultDocCommentUnmarshal, &gentenums.Gents.Stringers.All[0]
}

type GentEnumJsonMethods struct {
	gent.Opts

	DocCommentMarshal   string
	DocCommentUnmarshal string
	StringerToUse       *gentenums.StringMethodOpts
}

func (this *GentEnumJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.T.Method("MarshalJSON").Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(this.DocCommentMarshal).
		Code(
			ˇ.R.Set(T.SliceOf.Bytes.From(ctx.Import("strconv").C("Quote", This.C(this.StringerToUse.Name)))),
		)
}

func (this *GentEnumJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method("UnmarshalJSON", ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(this.DocCommentUnmarshal).
		Code(
			Var(ˇ.S.Name, T.String, nil),
			Tup(ˇ.S, ˇ.Err).Set(ctx.Import("strconv").C("Unquote", T.String.From(ˇ.V))),
			If(ˇ.Err.Eq(B.Nil), Then(
				Var(ˇ.T.Name, t.G.T, nil),
				Tup(ˇ.T, ˇ.Err).Set(N(this.StringerToUse.Parser.FuncName.With("T", t.Name, "str", this.StringerToUse.Name)).Of(ˇ.S)),
				If(ˇ.Err.Eq(B.Nil), Then(
					This.Deref().Set(ˇ.T),
				)),
			)),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentEnumJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsEnumish() {
		yield.Add(
			this.genMarshalMethod(ctx, t),
			this.genUnmarshalMethod(ctx, t),
		)
	}
	return
}
