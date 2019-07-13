package gentjson

import (
	"strconv"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.Structs.Marshal.Name, Gents.Structs.Marshal.DocComment, Gents.Structs.Marshal.InitialBytesCap =
		DefaultMethodNameMarshal, DefaultDocCommentMarshal, 128
	Gents.Structs.Unmarshal.Name, Gents.Structs.Unmarshal.DocComment =
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
		if !me.Marshal.Disabled {
			yield.Add(me.genMarshalMethod(ctx, t))
		}
		if !me.Unmarshal.Disabled {
			yield.Add(me.genUnmarshalMethod(ctx, t))
		}
	}
	return
}

func (me *GentStructJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) (ret *SynFunc) {
	ret = t.G.Tª.Method(me.Marshal.Name).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment.With("N", me.Marshal.Name)).
		Code(
			ˇ.R.Set(B.Make.Of(T.SliceOf.Bytes, 0, me.Marshal.InitialBytesCap)),
		)
	ret.Body.Add(me.genMarshalStruct(ctx, func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef }, nil)...)
	return
}

func (me *GentStructJsonMethods) genMarshalBasedOnType(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) (code Syns) {
	acc, t := field()
	switch {
	case t.ArrOrSlice.Of != nil:
		code.Add(me.genMarshalArrayOrSlice(ctx, field, writeName, jsonOmitEmpty))
	case t.Map.OfKey != nil:
		code.Add(writeName, ˇ.R.Set(B.Append.Of(ˇ.R, "{}").Spreads()))
	case t.Struct != nil:
		code.Add(me.genMarshalStruct(ctx, field, writeName)...)
	case t.Pointer.Of != nil:
		code.Add(me.genMarshalPointer(ctx, field, writeName, jsonOmitEmpty))
	case t.IsBuiltinPrimType(false):
		code.Add(me.genMarshalBuiltinPrim(ctx, field, writeName, jsonOmitEmpty, false))
	case t.Interface != nil:
		code.Add(me.genMarshalUnknown(ctx, field, writeName, jsonOmitEmpty, "", false))
	case t.Named.TypeName != "" && t.Named.PkgName != "":
		code.Add(me.genMarshalUnknown(ctx, field, writeName, jsonOmitEmpty, "", false))
	case t.Named.TypeName != "" && t.Named.PkgName == "":
		var mname string
		if gt := ctx.Pkg.Types.Named(t.Named.TypeName); gt != nil {
			isenumish := gt.IsEnumish()
			_ = ctx.GentExistsFor(gt, func(g gent.IGent) (ok bool) {
				gje, ok1 := g.(*GentEnumJsonMethods)
				gjs, ok2 := g.(*GentStructJsonMethods)
				if ok = ok1 && !gje.Disabled; ok {
					mname, isenumish = gje.Marshal.Name, true
				} else if ok = ok2 && !gjs.Disabled; ok {
					mname = gjs.Marshal.Name
				}
				return
			})
			if isenumish && mname == "" {
				code.Add(me.genMarshalBuiltinPrim(ctx, field, writeName, jsonOmitEmpty, true))
			} else if mname != "" {
				code.Add(me.genMarshalUnknown(ctx, field, writeName, jsonOmitEmpty, mname, isenumish))
			} else {
				code.Add(me.genMarshalBasedOnType(ctx, func() (ISyn, *TypeRef) {
					return acc, gt.Expr.GenRef
				}, writeName, jsonOmitEmpty)...)
			}
		} else {
			panic(t.Named.TypeName)
		}
	default:
		code.Add(writeName, ˇ.R.Set(B.Append.Of(ˇ.R, "?null").Spreads()))
	}
	return
}

func (me *GentStructJsonMethods) genMarshalStruct(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn) (code Syns) {
	acc, t := field()
	if writeName != nil {
		code.Add(writeName)
	}
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '{')))
	for i := range t.Struct.Fields {
		fld := &t.Struct.Fields[i]
		if jsonfieldname := fld.JsonNameFinal(); jsonfieldname != "" {
			pref, jsonomitempty := ",", fld.JsonOmitEmpty()
			if i == 0 {
				pref = ""
			}
			writename := ˇ.R.Set(B.Append.Of(ˇ.R, pref+strconv.Quote(jsonfieldname)+":").Spreads())
			code.Add(me.genMarshalBasedOnType(ctx,
				func() (ISyn, *TypeRef) { return D(acc, N(fld.EffectiveName())), fld.Type },
				writename, jsonomitempty)...)
		}
	}
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '}')))
	return
}

func (me *GentStructJsonMethods) genMarshalPointer(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ft := field()
	return If(B.Nil.Neq(facc), me.genMarshalBasedOnType(ctx, func() (ISyn, *TypeRef) {
		return facc, ft.Pointer.Of
	}, writeName, false),
		L(!jsonOmitEmpty), Then(
			writeName,
			ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()),
		))
}

func (me *GentStructJsonMethods) genMarshalArrayOrSlice(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ftype := field()
	iter, hasval := ctx.N(), ftype.IsntZeroish(facc, true, false)
	facci := func() (ISyn, *TypeRef) { return At(facc, iter), ftype.ArrOrSlice.Of }
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, '[')),
		ForEach(iter, None, facc,
			me.genMarshalBasedOnType(ctx,
				facci,
				If(iter.Neq(0), ˇ.R.Set(B.Append.Of(ˇ.R, ','))),
				false)...,
		),
		ˇ.R.Set(B.Append.Of(ˇ.R, ']')),
	))
}

func (*GentStructJsonMethods) genMarshalBuiltinPrim(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool, forceInt bool) ISyn {
	facc, ftype := field()
	pkgstrconv, hasval := ctx.Import("strconv"), ftype.IsntZeroish(facc, false, forceInt)
	var writeval ISyn
	if forceInt {
		writeval = pkgstrconv.C("FormatInt", T.Int64.From(facc), 10)
	} else {
		switch ftn := ftype.Named.TypeName; ftn {
		case T.Bool.Named.TypeName:
			writeval = pkgstrconv.C("FormatBool", facc)
		case T.Byte.Named.TypeName, T.Uint.Named.TypeName, T.Uint16.Named.TypeName, T.Uint32.Named.TypeName, T.Uint64.Named.TypeName, T.Uint8.Named.TypeName:
			writeval = pkgstrconv.C("FormatUint", T.Uint64.From(facc), 10)
		case T.Int.Named.TypeName, T.Int16.Named.TypeName, T.Int32.Named.TypeName, T.Int64.Named.TypeName, T.Int8.Named.TypeName:
			writeval = pkgstrconv.C("FormatInt", T.Int64.From(facc), 10)
		case T.Float32.Named.TypeName:
			writeval = pkgstrconv.C("FormatFloat", T.Float64.From(facc), 'f', -1, 32)
		case T.Float64.Named.TypeName:
			writeval = pkgstrconv.C("FormatFloat", facc, 'f', -1, 64)
		case T.Rune.Named.TypeName:
			writeval = pkgstrconv.C("Quote", T.String.From(facc))
		case T.String.Named.TypeName:
			writeval = pkgstrconv.C("Quote", facc)
		default:
			panic(ftn)
		}
	}
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, writeval).Spreads()),
	))
}

func (*GentStructJsonMethods) genMarshalUnknown(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool, implMethodName string, isKnownNumeric bool) ISyn {
	facc, ftype := field()
	pkgjson, canimpl := ctx.Import("encoding/json"), ftype.Named.PkgName == "" && (!ftype.CanNeverImplement()) && !isKnownNumeric
	hasval := ftype.IsntZeroish(facc, false, isKnownNumeric)
	ifnotnil := If(L(isKnownNumeric && !jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		Var(ˇ.E.Name, T.Error, nil),
		Var(ˇ.Sl.Name, T.SliceOf.Bytes, nil),
		GEN_IF(implMethodName != "", Then(
			Tup(ˇ.Sl, ˇ.E).Set(Call(D(facc, N(implMethodName)))),
		), Else(
			GEN_IF(canimpl, Then(
				Tup(ˇ.J, ˇ.Ok).Let(D(facc, pkgjson.T("Marshaler"))),
				If(ˇ.Ok.And(ˇ.J.Neq(B.Nil)), Then(
					Tup(ˇ.Sl, ˇ.E).Set(ˇ.J.C("MarshalJSON")),
				), Else(
					Tup(ˇ.Sl, ˇ.E).Set(pkgjson.C("Marshal", facc)),
				)),
			), Else(
				Tup(ˇ.Sl, ˇ.E).Set(pkgjson.C("Marshal", facc)),
			)),
		)),
		If(ˇ.E.Eq(B.Nil), Then(
			ˇ.R.Set(B.Append.Of(ˇ.R, ˇ.Sl).Spreads()),
		), Else(
			ˇ.Err.Set(ˇ.E),
			Ret(nil),
		)),
	))
	if (!jsonOmitEmpty) && !isKnownNumeric {
		ifnotnil.Else.Add(writeName, ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()))
	}
	return ifnotnil
}

func (me *GentStructJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code()
}
