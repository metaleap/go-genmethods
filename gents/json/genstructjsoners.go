package gentjson

import (
	"strconv"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.Structs.Marshal.MethodName, Gents.Structs.Marshal.DocComment, Gents.Structs.Marshal.InitialBytesCap =
		DefaultMethodNameMarshal, DefaultDocCommentMarshal, 128
	Gents.Structs.Unmarshal.MethodName, Gents.Structs.Unmarshal.DocComment =
		DefaultMethodNameUnmarshal, DefaultDocCommentUnmarshal
}

type GentStructJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
		InitialBytesCap int
	}
	Unmarshal struct {
		JsonMethodOpts
	}
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

func (me *GentStructJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) (ret *SynFunc) {
	ret = t.G.Tª.Method(me.Marshal.MethodName).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment).
		Code(
			ˇ.R.Set(B.Make.Of(T.SliceOf.Bytes, 1, me.Marshal.InitialBytesCap)),
			ˇ.R.At(0).Set('{'),
		)
	ts, code := t.Expr.GenRef.Struct, &ret.Body
	for i := range ts.Fields {
		tsf := &ts.Fields[i]
		if jsonfieldname := tsf.JsonNameFinal(); jsonfieldname != "" {
			pref, jsonomitempty := ",", tsf.JsonOmitEmpty()
			if i == 0 {
				pref = ""
			}
			writename := ˇ.R.Set(B.Append.Of(ˇ.R, pref+strconv.Quote(jsonfieldname)+":").Spreads())
			switch tsft, tsfn, pkgjson := tsf.Type, tsf.Name, ctx.Import("encoding/json"); true {
			case tsft.ArrOrSlice.Of != nil:
				code.Add(writename, ˇ.R.Set(B.Append.Of(ˇ.R, "[]").Spreads()))
			case tsft.Map.OfKey != nil:
				code.Add(writename, ˇ.R.Set(B.Append.Of(ˇ.R, "{}").Spreads()))
			case tsft.Interface != nil:
				code.Add(me.genMarshalInterface(writename, jsonomitempty, pkgjson, tsfn))
			default:
				code.Add(writename, ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()))
			}
		}
	}
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '}')))
	return
}

func (*GentStructJsonMethods) genMarshalInterface(writeName ISyn, jsonOmitEmpty bool, pkgJson PkgName, fieldName string) ISyn {
	ifnotnil := If(Self.D(fieldName).Neq(B.Nil), Then(
		writeName,
		Var(ˇ.E.Name, T.Error, nil),
		Var(ˇ.Sl.Name, T.SliceOf.Bytes, nil),
		Tup(ˇ.J, ˇ.Ok).Let(Self.D(fieldName).D(pkgJson.T("Marshaler"))),
		If(ˇ.Ok.And(ˇ.J.Neq(B.Nil)), Then(
			Tup(ˇ.Sl, ˇ.E).Set(ˇ.J.C("MarshalJSON")),
		), Else(
			Tup(ˇ.Sl, ˇ.E).Set(pkgJson.C("Marshal", Self.D(fieldName))),
		)),
		If(ˇ.E.Eq(B.Nil), Then(
			ˇ.R.Set(B.Append.Of(ˇ.R, ˇ.Sl).Spreads()),
		), Else(
			ˇ.Err.Set(ˇ.E),
			Ret(nil),
		)),
	))
	if !jsonOmitEmpty {
		ifnotnil.Else.Add(writeName, ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()))
	}
	return ifnotnil
}

func (me *GentStructJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.MethodName, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment).
		Code()
}
