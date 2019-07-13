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
			switch tsft := tsf.Type; true {
			case tsft.ArrOrSlice.Of != nil:
				code.Add(writename, ˇ.R.Set(B.Append.Of(ˇ.R, "[]").Spreads()))
			case tsft.Map.OfKey != nil:
				code.Add(writename, ˇ.R.Set(B.Append.Of(ˇ.R, "{}").Spreads()))
			case tsft.Interface != nil:
				code.Add(me.genMarshalInterface(ctx, tsf, writename, jsonomitempty))
			case tsft.IsBuiltinPrimType(false):
				code.Add(me.genMarshalBuiltinPrim(ctx, tsf, writename, jsonomitempty))
			default:
				code.Add(writename, ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()))
			}
		}
	}
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '}')))
	return
}

func (*GentStructJsonMethods) genMarshalBuiltinPrim(ctx *gent.Ctx, field *SynStructField, writeName ISyn, jsonOmitEmpty bool) ISyn {
	fieldacc, pkgstrconv := Self.D(field.Name), ctx.Import("strconv")
	hasval := field.Type.IsntZeroish(fieldacc, false, false)
	var writeval ISyn
	switch ftn := field.Type.Named.TypeName; ftn {
	case T.Bool.Named.TypeName:
		writeval = pkgstrconv.C("FormatBool", fieldacc)
	case T.Byte.Named.TypeName, T.Uint.Named.TypeName, T.Uint16.Named.TypeName, T.Uint32.Named.TypeName, T.Uint64.Named.TypeName, T.Uint8.Named.TypeName:
		writeval = pkgstrconv.C("FormatUint", T.Uint64.From(fieldacc), 10)
	case T.Int.Named.TypeName, T.Int16.Named.TypeName, T.Int32.Named.TypeName, T.Int64.Named.TypeName, T.Int8.Named.TypeName:
		writeval = pkgstrconv.C("FormatInt", T.Int64.From(fieldacc), 10)
	case T.Float32.Named.TypeName:
		writeval = pkgstrconv.C("FormatFloat", T.Float64.From(fieldacc), 'f', -1, 32)
	case T.Float64.Named.TypeName:
		writeval = pkgstrconv.C("FormatFloat", fieldacc, 'f', -1, 64)
	case T.Rune.Named.TypeName: // not sure if handled already in int32 case above, not pressing right now
		writeval = pkgstrconv.C("Quote", T.String.From(fieldacc))
	case T.String.Named.TypeName:
		writeval = pkgstrconv.C("Quote", fieldacc)
	default:
		panic(ftn)
	}
	writeempty := !jsonOmitEmpty
	return If(L(writeempty).Or(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, writeval).Spreads()),
	))
}

func (*GentStructJsonMethods) genMarshalInterface(ctx *gent.Ctx, field *SynStructField, writeName ISyn, jsonOmitEmpty bool) ISyn {
	fieldacc, pkgjson := Self.D(field.Name), ctx.Import("encoding/json")
	ifnotnil := If(fieldacc.Neq(B.Nil), Then(
		writeName,
		Var(ˇ.E.Name, T.Error, nil),
		Var(ˇ.Sl.Name, T.SliceOf.Bytes, nil),
		Tup(ˇ.J, ˇ.Ok).Let(fieldacc.D(pkgjson.T("Marshaler"))),
		If(ˇ.Ok.And(ˇ.J.Neq(B.Nil)), Then(
			Tup(ˇ.Sl, ˇ.E).Set(ˇ.J.C("MarshalJSON")),
		), Else(
			Tup(ˇ.Sl, ˇ.E).Set(pkgjson.C("Marshal", fieldacc)),
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
