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
		if pref, jsonfieldname := ",", tsf.JsonNameFinal(); jsonfieldname != "" {
			if i == 0 {
				pref = ""
			}
			code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, pref+strconv.Quote(jsonfieldname)+":").Spreads()))
			switch tsft := tsf.Type; true {
			case tsft.ArrOrSlice.Of != nil:
				code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, "[]").Spreads()))
			case tsft.Map.OfKey != nil:
				code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, "{}").Spreads()))
			case tsft.Interface != nil:
				/*
					if jm, ok := me.Val.(json.Marshaler); ok && jm != nil {
						if b, e := jm.MarshalJSON(); e != nil {
							err = e
							return
						} else {
							r = append(r, b...)
						}
					}
				*/
				code.Add(Block(
					Tup(ˇ.J, ˇ.Ok).Let(Self.D("Val").D(TFrom(ctx.Import("encoding/json"), "Marshaler"))),
					If(ˇ.Ok.And(ˇ.J.Neq(B.Nil)), Then(
						Tup(ˇ.Sl, ˇ.E).Let(ˇ.J.C("MarshalJSON")),
						If(ˇ.E.Eq(B.Nil), Then(
							ˇ.R.Set(B.Append.Of(ˇ.R, ˇ.Sl).Spreads()),
						), Else(
							ˇ.Err.Set(ˇ.E),
							Ret(nil),
						)),
					)),
				))
			default:
				code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()))
			}
		}
	}
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '}')))
	return
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
