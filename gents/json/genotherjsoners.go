package gentjson

import (
	"strconv"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.OtherTypes.Marshal.Name, Gents.OtherTypes.Marshal.DocComment, Gents.OtherTypes.Marshal.InitialBytesCap =
		DefaultMethodNameMarshal, DefaultDocCommentMarshal, 64
	Gents.OtherTypes.Unmarshal.Name, Gents.OtherTypes.Unmarshal.DocComment =
		DefaultMethodNameUnmarshal, DefaultDocCommentUnmarshal
}

type GentTypeJsonMethods struct {
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
func (me *GentTypeJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if !t.IsEnumish() {
		if (!me.Marshal.Disabled) && (me.Marshal.MayGenFor == nil || me.Marshal.MayGenFor(t)) {
			yield.Add(me.genMarshalMethod(ctx, t))
		}
		if (!me.Unmarshal.Disabled) && (me.Unmarshal.MayGenFor == nil || me.Unmarshal.MayGenFor(t)) {
			yield.Add(me.genUnmarshalMethod(ctx, t))
		}
	}
	return
}

func (me *GentTypeJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	self := t.G.Tª
	var kickoff ISyn
	if t.Expr.GenRef.Struct != nil {
		kickoff = me.genMarshalStruct(ctx, func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef }, nil)
	} else {
		self = t.G.T
		kickoff = me.genMarshalBasedOnType(ctx, func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef }, nil, false)
	}
	return self.Method(me.Marshal.Name).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment.With("N", me.Marshal.Name)).
		Code(
			ˇ.R.Set(B.Make.Of(T.SliceOf.Bytes, 0, me.Marshal.InitialBytesCap)),
			If(Self.Eq(B.Nil), Then(
				ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()),
			), Else(
				kickoff,
			)),
		)
}

func (me *GentTypeJsonMethods) genMarshalBasedOnType(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) (code Syns) {
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
				gjs, ok2 := g.(*GentTypeJsonMethods)
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

func (me *GentTypeJsonMethods) genMarshalStruct(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn) (code Syns) {
	acc, t := field()
	if writeName != nil {
		code.Add(writeName)
	}
	idx := ctx.N("si")
	code.Add(
		ˇ.R.Set(B.Append.Of(ˇ.R, '{')),
		idx.Let(B.Len.Of(ˇ.R)),
	)
	code.Add(me.genMarshalStructFields(ctx, t.Struct.Fields, acc)...)
	code.Add(
		ˇ.R.Set(B.Append.Of(ˇ.R, '}')),
		If(ˇ.R.At(idx).Eq(','), Then(ˇ.R.At(idx).Set(' '))),
	)
	return
}

func (me *GentTypeJsonMethods) genMarshalStructFields(ctx *gent.Ctx, fields SynStructFields, acc ISyn) (code Syns) {
	for i := range fields {
		fld := &fields[i]
		var fldcode Syns
		if ft := ctx.Pkg.Types.Named(fld.Type.UltimateElemType().Named.TypeName); fld.Name == "" && ft != nil && ft.Expr.GenRef.Struct != nil {
			fldcode = me.genMarshalStructFields(ctx, ft.Expr.GenRef.Struct.Fields, D(acc, N(fld.EffectiveName())))
		} else if jsonfieldname := fld.JsonNameFinal(); jsonfieldname != "" {
			jsonomitempty := fld.JsonOmitEmpty()
			writename := ˇ.R.Set(B.Append.Of(ˇ.R, ","+strconv.Quote(jsonfieldname)+":").Spreads())
			fldcode = me.genMarshalBasedOnType(ctx,
				func() (ISyn, *TypeRef) { return D(acc, N(fld.EffectiveName())), fld.Type },
				writename, jsonomitempty)
		}
		code.Add(fldcode...)
	}
	return
}

func (me *GentTypeJsonMethods) genMarshalPointer(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
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

func (me *GentTypeJsonMethods) genMarshalArrayOrSlice(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ftype := field()
	iter, idx, hasval := ctx.N("i"), ctx.N("ai"), ftype.IsntZeroish(facc, true, false)
	facci := func() (ISyn, *TypeRef) { return At(facc, iter), ftype.ArrOrSlice.Of }
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, '[')),
		idx.Let(B.Len.Of(ˇ.R)),
		ForEach(iter, None, facc,
			me.genMarshalBasedOnType(ctx,
				facci,
				ˇ.R.Set(B.Append.Of(ˇ.R, ',')),
				false)...,
		),
		ˇ.R.Set(B.Append.Of(ˇ.R, ']')),
		If(ˇ.R.At(idx).Eq(','), Then(ˇ.R.At(idx).Set(' '))),
	))
}

func (me *GentTypeJsonMethods) genMarshalMap(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool) ISyn {
	facc, ftype := field()
	hasval, idx, iterk, iterv := ftype.IsntZeroish(facc, true, false), ctx.N("mi"), ctx.N("mk"), ctx.N("mv")
	fkacc, fvacc := func() (ISyn, *TypeRef) { return iterk, ftype.Map.OfKey }, func() (ISyn, *TypeRef) { return iterv, ftype.Map.ToVal }
	key2str := me.genToString(ctx, fkacc, false, true)
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, '{')),
		idx.Let(B.Len.Of(ˇ.R)),
		ForEach(iterk, iterv, facc,
			me.genMarshalBasedOnType(ctx,
				fvacc,
				Block(
					ˇ.R.Set(B.Append.Of(ˇ.R, ',')),
					ˇ.R.Set(B.Append.Of(ˇ.R, key2str).Spreads()),
					ˇ.R.Set(B.Append.Of(ˇ.R, ':')),
				),
				false)...,
		),
		ˇ.R.Set(B.Append.Of(ˇ.R, '}')),
		If(ˇ.R.At(idx).Eq(','), Then(ˇ.R.At(idx).Set(' '))),
	))
}

func (*GentTypeJsonMethods) genToString(ctx *gent.Ctx, field func() (ISyn, *TypeRef), forceInt bool, ensureQuoted bool) (code ISyn) {
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

func (me *GentTypeJsonMethods) genMarshalBuiltinPrim(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool, forceInt bool) ISyn {
	facc, ftype := field()
	hasval := ftype.IsntZeroish(facc, false, forceInt)
	writeval := me.genToString(ctx, field, forceInt, false)
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, writeval).Spreads()),
	))
}

func (*GentTypeJsonMethods) genMarshalUnknown(ctx *gent.Ctx, field func() (ISyn, *TypeRef), writeName ISyn, jsonOmitEmpty bool, implMethodName string, isKnownNumeric bool) ISyn {
	facc, ftype := field()
	pkgjson, canimpl := ctx.Import("encoding/json"), ftype.Named.PkgName == "" && (!ftype.CanNeverImplement()) && !isKnownNumeric
	hasval := ftype.IsntZeroish(facc, false, isKnownNumeric)
	ifnotnil := If(L(implMethodName != "" || (isKnownNumeric && !jsonOmitEmpty)).Or1(hasval), Then(
		writeName,
		Var(ˇ.E.Name, T.Error, nil),
		Var(ˇ.Sl.Name, T.SliceOf.Bytes, nil),
		GEN_IF(implMethodName != "", Then(
			Tup(ˇ.Sl, ˇ.E).Set(Call(D(facc, N(implMethodName)))),
		), Else(
			GEN_IF(canimpl, Then(
				Tup(ˇ.J, Nope).Let(D(facc, pkgjson.T("Marshaler"))),
				If(ˇ.J.Neq(B.Nil), Then(
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

func (me *GentTypeJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code()
}
