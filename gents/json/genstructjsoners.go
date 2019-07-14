package gentjson

import (
	"strconv"

	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
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
		code.Add(me.genMarshalMap(ctx, field, writeName, jsonOmitEmpty))
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
		panic(t.String())
	}
	return
}

func (me *GentStructJsonMethods) genMarshalStruct(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn) (code Syns) {
	acc, t := field()
	if writeName != nil {
		code.Add(writeName)
	}
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '{')))
	code.Add(me.genMarshalStructFields(ctx, t.Struct.Fields, acc, true)...)
	code.Add(ˇ.R.Set(B.Append.Of(ˇ.R, '}')))
	return
}

func (me *GentStructJsonMethods) genMarshalStructFields(ctx *gent.Ctx, fields SynStructFields, acc ISyn, skipLeadingComma bool) (code Syns) {
	for i := range fields {
		fld := &fields[i]
		if ft := ctx.Pkg.Types.Named(fld.Type.UltimateElemType().Named.TypeName); fld.Name == "" && ft != nil && ft.Expr.GenRef.Struct != nil {
			code.Add(me.genMarshalStructFields(ctx, ft.Expr.GenRef.Struct.Fields, D(acc, N(fld.EffectiveName())), skipLeadingComma && i == 0)...)
		} else if jsonfieldname := fld.JsonNameFinal(); jsonfieldname != "" {
			jsonomitempty := fld.JsonOmitEmpty()
			writename := ˇ.R.Set(B.Append.Of(ˇ.R, ustr.If(skipLeadingComma && i == 0, "", ",")+strconv.Quote(jsonfieldname)+":").Spreads())
			code.Add(me.genMarshalBasedOnType(ctx,
				func() (ISyn, *TypeRef) { return D(acc, N(fld.EffectiveName())), fld.Type },
				writename, jsonomitempty)...)
		}
	}
	return
}

func (me *GentStructJsonMethods) genMarshalPointer(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ft := field()
	return If(
		B.Nil.Neq(facc), Then(
			me.genMarshalBasedOnType(ctx, func() (ISyn, *TypeRef) {
				return facc, ft.Pointer.Of
			}, writeName, false),
		),
		L(!jsonOmitEmpty), Then(
			writeName,
			ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()),
		))
}

func (me *GentStructJsonMethods) genMarshalArrayOrSlice(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ftype := field()
	iter, hasval := ctx.N("i"), ftype.IsntZeroish(facc, true, false)
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

func (me *GentStructJsonMethods) genMarshalMap(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ftype := field()
	hasval, iterk, iterv, isfirst := ftype.IsntZeroish(facc, true, false), ctx.N("mk"), ctx.N("mv"), ctx.N("mf")
	fkacc, fvacc := func() (ISyn, *TypeRef) { return iterk, ftype.Map.OfKey }, func() (ISyn, *TypeRef) { return iterv, ftype.Map.ToVal }
	key2str := me.genToString(ctx, fkacc, false, true)
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		isfirst.Let(true),
		ˇ.R.Set(B.Append.Of(ˇ.R, '{')),
		ForEach(iterk, iterv, facc,
			me.genMarshalBasedOnType(ctx,
				fvacc,
				Block(
					If(isfirst, Then(isfirst.Set(false)), Else(ˇ.R.Set(B.Append.Of(ˇ.R, ',')))),
					ˇ.R.Set(B.Append.Of(ˇ.R, key2str).Spreads()),
					ˇ.R.Set(B.Append.Of(ˇ.R, ':')),
				),
				false)...,
		),
		ˇ.R.Set(B.Append.Of(ˇ.R, '}')),
	))
}

func (*GentStructJsonMethods) genToString(ctx *gent.Ctx, field func() (ISyn, *TypeRef), forceInt bool, ensureQuoted bool) (code ISyn) {
	isquoted, pkgstrconv := false, ctx.Import("strconv")
	facc, ftype := field()
	if forceInt {
		code = pkgstrconv.C("FormatInt", T.Int64.From(facc), 10)
	} else {
		switch ftn := ftype.Named; ftn {
		case T.Bool.Named:
			code = pkgstrconv.C("FormatBool", facc)
		case T.Byte.Named, T.Uint.Named, T.Uint16.Named, T.Uint32.Named, T.Uint64.Named, T.Uint8.Named:
			code = pkgstrconv.C("FormatUint", T.Uint64.From(facc), 10)
		case T.Int.Named, T.Int16.Named, T.Int32.Named, T.Int64.Named, T.Int8.Named:
			code = pkgstrconv.C("FormatInt", T.Int64.From(facc), 10)
		case T.Float32.Named:
			code = pkgstrconv.C("FormatFloat", T.Float64.From(facc), 'f', -1, 32)
		case T.Float64.Named:
			code = pkgstrconv.C("FormatFloat", facc, 'f', -1, 64)
		case T.Rune.Named:
			code, isquoted = pkgstrconv.C("Quote", T.String.From(facc)), true
		case T.String.Named:
			code, isquoted = pkgstrconv.C("Quote", facc), true
		default:
			code = Call(D(facc, N("String")))
		}
	}
	if ensureQuoted && !isquoted {
		code = pkgstrconv.C("Quote", code)
	}
	return
}

func (me *GentStructJsonMethods) genMarshalBuiltinPrim(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool, forceInt bool) ISyn {
	facc, ftype := field()
	hasval := ftype.IsntZeroish(facc, false, forceInt)
	writeval := me.genToString(ctx, field, forceInt, false)
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
