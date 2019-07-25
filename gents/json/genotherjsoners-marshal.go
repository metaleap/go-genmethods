package gentjson

import (
	"strconv"

	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
)

func (me *GentTypeJsonMethods) genMarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var selftype *TypeRef
	var code ISyn
	if t.Expr.GenRef.Struct != nil {
		selftype, code = t.G.Tª, If(Self.Eq(B.Nil), Then(
			jsonWriteNull,
		), Else(
			me.genMarshalStruct(ctx, Self, t.Expr.GenRef, nil),
		))
	} else {
		selftype, code = t.G.T, me.genMarshalBasedOnType(ctx, Self, t.Expr.GenRef, nil, false, true)
	}
	return selftype.Method(me.Marshal.Name).Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
		Doc(me.Marshal.DocComment.With("N", me.Marshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				B.Panic.Of(t.Name),
			), Else(
				ˇ.R.Set(B.Make.Of(T.SliceOf.Bytes, 0, me.Marshal.InitialBytesCap)),
				code,
			)),
		)
}

func (me *GentTypeJsonMethods) genMarshalBasedOnType(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn, jsonOmitEmpty bool, canUseExtraDef bool) (code Syns) {
	me.genExtraDefs(ctx)

	if canUseExtraDef && !fType.IsBuiltinPrimType(false) {
		for _, tref := range me.Marshal.TryInterfaceTypesBeforeStdlib {
			if tref.Equiv(fType) {
				defname := me.genExtraDefName(ctx, fType)
				var writefromdefcall ISyn = Block(
					Tup(ˇ.Sl, ˇ.E).Let(C(defname, fAcc)),
					If(ˇ.E.Eq(B.Nil), Then(
						writeName,
						ˇ.R.Set(B.Append.Of(ˇ.R, ˇ.Sl).Spreads()),
					), Else(
						ˇ.Err.Set(ˇ.E),
						Ret(nil),
					)),
				)
				if !jsonOmitEmpty {
					code.Add(writefromdefcall)
				} else {
					code.Add(If(fType.IsntZeroish(fAcc, false, false),
						writefromdefcall))
				}
				return
			}
		}
	}

	switch {
	case fType.ArrOrSlice.Of != nil:
		code.Add(me.genMarshalArrayOrSlice(ctx, fAcc, fType, writeName, jsonOmitEmpty))
	case fType.Map.OfKey != nil:
		code.Add(me.genMarshalMap(ctx, fAcc, fType, writeName, jsonOmitEmpty))
	case fType.Struct != nil:
		code.Add(me.genMarshalStruct(ctx, fAcc, fType, writeName)...)
	case fType.Pointer.Of != nil:
		code.Add(me.genMarshalPointer(ctx, fAcc, fType, writeName, jsonOmitEmpty))
	case fType.IsBuiltinPrimType(false):
		code.Add(me.genMarshalBuiltinPrim(ctx, fAcc, fType, writeName, jsonOmitEmpty, false))
	case fType.Interface != nil:
		code.Add(me.genMarshalUnknown(ctx, fAcc, fType, writeName, jsonOmitEmpty, "", false, false))
	case fType.Named.TypeName != "":
		var pkg *gent.Pkg
		if fType.Named.PkgName == "" {
			pkg = ctx.Pkg
		} else { // if pkg = gent.TryExtPkg(t.Named.PkgName); pkg == nil  /* ext-pkgs stuff not really working just yet, TODO when it becomes more pressing */ {
			code.Add(me.genMarshalUnknown(ctx, fAcc, fType, writeName, jsonOmitEmpty, "", false, false))
			return
		}

		var mname string
		if gt := pkg.Types.Named(fType.Named.TypeName); gt == nil {
			panic(fType.Named.TypeName)
		} else {
			isenumish := gt.IsEnumish()
			_ = ctx.GentExistsFor(gt, func(g gent.IGent) (ok bool) {
				gje, ok1 := g.(*GentEnumJsonMethods)
				gjt, ok2 := g.(*GentTypeJsonMethods)
				if ok = ok1 && !(gje.Disabled || gje.Marshal.Disabled); ok {
					mname, isenumish = gje.Marshal.Name, true
				} else if ok = ok2 && !(gjt.Disabled || gjt.Marshal.Disabled); ok {
					mname = gjt.Marshal.Name
				}
				return
			})
			if isenumish && mname == "" {
				code.Add(me.genMarshalBuiltinPrim(ctx, fAcc, fType, writeName, jsonOmitEmpty, true))
			} else if mname != "" {
				code.Add(me.genMarshalUnknown(ctx, fAcc, fType, writeName, jsonOmitEmpty,
					mname, isenumish || 0 != gt.Expr.GenRef.BitSizeIfBuiltInNumberType(),
					gt.IsSliceOrArray() || gt.Expr.GenRef.Map.OfKey != nil))
			} else {
				code.Add(me.genMarshalBasedOnType(ctx, fAcc, gt.Expr.GenRef, writeName, jsonOmitEmpty, true)...)
			}
		}
	default:
		panic(fType.String())
	}
	return
}

func (me *GentTypeJsonMethods) genMarshalStruct(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn) (code Syns) {
	if writeName != nil {
		code.Add(writeName)
	}
	idx := ctx.N("si")
	code.Add(
		ˇ.R.Set(B.Append.Of(ˇ.R, '{')),
		idx.Let(B.Len.Of(ˇ.R)),
	)
	code.Add(me.genMarshalStructFields(ctx, fAcc, fType)...)
	code.Add(
		ˇ.R.Set(B.Append.Of(ˇ.R, '}')),
		me.genResliceOrFixup(idx),
	)
	return
}

func (me *GentTypeJsonMethods) genMarshalStructFields(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (code Syns) {
	for i := range fType.Struct.Fields {
		fld := &fType.Struct.Fields[i]
		var fldcode Syns
		if ft := ctx.Pkg.Types.Named(fld.Type.UltimateElemType().Named.TypeName); fld.Name == "" && ft != nil && ft.Expr.GenRef.Struct != nil {
			fldcode = me.genMarshalStructFields(ctx, D(fAcc, N(fld.EffectiveName())), ft.Expr.GenRef)
		} else if jsonfieldname := fld.JsonNameFinal(); jsonfieldname != "" {
			jsonomitempty := fld.JsonOmitEmpty()
			writename := ˇ.R.Set(B.Append.Of(ˇ.R, ","+strconv.Quote(jsonfieldname)+":").Spreads())
			fldcode = me.genMarshalBasedOnType(ctx, D(fAcc, N(fld.EffectiveName())), fld.Type, writename, jsonomitempty, true)
		}
		code.Add(fldcode...)
	}
	return
}

func (me *GentTypeJsonMethods) genMarshalPointer(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn, jsonOmitEmpty bool) ISyn {
	return If(
		B.Nil.Neq(fAcc), Then(
			me.genMarshalBasedOnType(ctx, fAcc, fType.Pointer.Of, writeName, false, true),
		),
		L(!jsonOmitEmpty), Then(
			writeName,
			jsonWriteNull,
		))
}

func (me *GentTypeJsonMethods) genMarshalArrayOrSlice(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn, jsonOmitEmpty bool) ISyn {
	iter, idx, hasval := ctx.N("i"), ctx.N("ai"), fType.IsntZeroish(fAcc, true, false)

	writeiter := Block(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, '[')),
		idx.Let(B.Len.Of(ˇ.R)),
		ForEach(iter, None, fAcc,
			me.genMarshalBasedOnType(ctx, At(fAcc, iter), fType.ArrOrSlice.Of, ˇ.R.Set(B.Append.Of(ˇ.R, ',')), false, true)...,
		),
		ˇ.R.Set(B.Append.Of(ˇ.R, ']')),
		me.genResliceOrFixup(idx),
	)

	if jsonOmitEmpty {
		return If(hasval, writeiter)
	} else {
		return If(B.Nil.Eq(fAcc), Then(
			writeName,
			jsonWriteNull,
		),
			writeiter,
		)
	}
}

func (me *GentTypeJsonMethods) genMarshalMap(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn, jsonOmitEmpty bool) ISyn {
	hasval, idx, iterk, iterv := fType.IsntZeroish(fAcc, true, false), ctx.N("mi"), ctx.N("mk"), ctx.N("mv")
	key2str := me.genToString(ctx, iterk, fType.Map.OfKey, false, true)

	writeiter := Block(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, '{')),
		idx.Let(B.Len.Of(ˇ.R)),
		ForEach(iterk, iterv, fAcc,
			me.genMarshalBasedOnType(ctx, iterv, fType.Map.ToVal, Block(
				ˇ.R.Set(B.Append.Of(ˇ.R, ',')),
				ˇ.R.Set(B.Append.Of(ˇ.R, key2str).Spreads()),
				ˇ.R.Set(B.Append.Of(ˇ.R, ':')),
			), false, true)...,
		),
		ˇ.R.Set(B.Append.Of(ˇ.R, '}')),
		me.genResliceOrFixup(idx),
	)

	if jsonOmitEmpty {
		return If(hasval, writeiter)
	}

	return If(B.Nil.Eq(fAcc), Then(
		writeName, jsonWriteNull,
	),
		writeiter,
	)
}

func (me *GentTypeJsonMethods) genMarshalBuiltinPrim(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn, jsonOmitEmpty bool, forceInt bool) ISyn {
	hasval := fType.IsntZeroish(fAcc, false, forceInt)
	writeval := me.genToString(ctx, fAcc, fType, forceInt, false)
	return If(L(!jsonOmitEmpty).Or1(hasval), Then(
		writeName,
		ˇ.R.Set(B.Append.Of(ˇ.R, writeval).Spreads()),
	))
}

func (me *GentTypeJsonMethods) genMarshalUnknown(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, writeName ISyn, jsonOmitEmpty bool, implMethodName string, isKnownNumish bool, isKnownLenish bool) ISyn {
	pkgjson, canimpl := me.pkgjson, fType.Named.PkgName == "" && (!fType.CanNeverImplement()) && !(isKnownNumish || isKnownLenish)
	hasval := fType.IsntZeroish(fAcc, isKnownLenish, isKnownNumish)
	skipcheck := (!jsonOmitEmpty) || (implMethodName != "" && !isKnownLenish)
	stdlibfallback := me.Marshal.OnStdlibFallbacks(ctx, fAcc,
		Tup(ˇ.Sl, ˇ.E).Set(pkgjson.C("Marshal", fAcc)),
	)
	ifnotnil := If(L(skipcheck).Or1(hasval), Then(
		writeName,
		Var(ˇ.E.Name, T.Error, nil),
		Var(ˇ.Sl.Name, T.SliceOf.Bytes, nil),
		GEN_IF(implMethodName != "", Then(
			Tup(ˇ.Sl, ˇ.E).Set(Call(D(fAcc, N(implMethodName)))),
		), Else(
			GEN_IF(canimpl, Then(
				Tup(ˇ.J, ˇ.Ok).Let(D(fAcc, pkgjson.T("Marshaler"))),
				If(ˇ.Ok, Then(
					If(ˇ.J.Neq(B.Nil), Then(
						Tup(ˇ.Sl, ˇ.E).Set(ˇ.J.C("MarshalJSON")),
					), Else(
						ˇ.Sl.Set(T.SliceOf.Bytes.From("null")),
					)),
				), Else(
					me.genIfaceFallbacks(ctx, fAcc, stdlibfallback),
				)),
			), Else(
				stdlibfallback,
			)),
		)),
		If(ˇ.E.Eq(B.Nil), Then(
			ˇ.R.Set(B.Append.Of(ˇ.R, ˇ.Sl).Spreads()),
		), Else(
			ˇ.Err.Set(ˇ.E),
			Ret(nil),
		)),
	))
	if (!jsonOmitEmpty) && !isKnownNumish {
		ifnotnil.Else.Add(writeName, jsonWriteNull)
	}
	return ifnotnil
}

func (me *GentTypeJsonMethods) genExtraDefName(ctx *gent.Ctx, t *TypeRef) string {
	return ctx.Opt.HelpersPrefix + me.Marshal.HelpersPrefix + ustr.ReplB(t.String(), '[', 's', ']', '_', '*', 'p', '{', '_', '}', '_', '.', '_')
}

func (me *GentTypeJsonMethods) genExtraDefs(ctx *gent.Ctx) {
	if !me.Marshal.tryInterfaceTypesDefsDone {
		me.Marshal.tryInterfaceTypesDefsDone = true
		defsdone := map[string]struct{}{}
		for _, checktype := range me.Marshal.TryInterfaceTypesBeforeStdlib {
			defname := me.genExtraDefName(ctx, checktype)
			var defcode ISyn = me.genMarshalBasedOnType(ctx, ˇ.V, checktype, Block(), false, false)
			if checktype.Interface != nil {
				defcode = If(B.Nil.Eq(ˇ.V), Then(
					ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads()),
				), Else(
					defcode,
				))
			}
			if _, defdone := defsdone[defname]; !defdone {
				defsdone[defname] = struct{}{}
				ctx.ExtraDefs = append(ctx.ExtraDefs, Func(defname, ˇ.V.OfType(checktype)).
					Rets(ˇ.R.OfType(T.SliceOf.Bytes), ˇ.Err).
					Code(defcode),
				)
			}
		}
	}
}

func (me *GentTypeJsonMethods) genIfaceFallbacks(ctx *gent.Ctx, fAcc ISyn, stdlibFallback ISyn) (ifaceFallbacks ISyn) {
	ifaceFallbacks = stdlibFallback
	me.genExtraDefs(ctx)
	for _, checktype := range me.Marshal.TryInterfaceTypesBeforeStdlib {
		if !checktype.IsEmptyInterface() {
			defname := me.genExtraDefName(ctx, checktype)
			val, ok := ctx.N("v"), ctx.N("ok")
			ifaceFallbacks = Block(
				Tup(val, ok).Let(D(fAcc, checktype)),
				If(ok, Then(
					Tup(ˇ.Sl, ˇ.E).Set(C(defname, val)),
				),
					ifaceFallbacks),
			)
		}
	}

	return
}

func (me *GentTypeJsonMethods) genResliceOrFixup(idx Named) ISyn {
	if me.Marshal.ResliceInsteadOfWhitespace {
		return If(ˇ.R.At(idx).Eq(','), Then(
			ˇ.R.Set(B.Append.Of(ˇ.R.Sl(None, idx), ˇ.R.Sl(idx.Plus(1), None)).Spreads()),
		))
	}
	return If(ˇ.R.At(idx).Eq(','), Then(
		ˇ.R.At(idx).Set(' '),
	))
}

func (*GentTypeJsonMethods) genToString(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, forceInt bool, ensureQuoted bool) (code ISyn) {
	isquoted, pkgstrconv := false, ctx.Import("strconv")
	if forceInt {
		code = pkgstrconv.C("FormatInt", T.Int64.From(fAcc), 10)
	} else {
		switch ftn := fType.Named; ftn {
		case T.Bool.Named:
			code = pkgstrconv.C("FormatBool", fAcc)
		case T.Byte.Named, T.Uint.Named, T.Uint16.Named, T.Uint32.Named, T.Uint64.Named, T.Uint8.Named:
			code = pkgstrconv.C("FormatUint", T.Uint64.From(fAcc), 10)
		case T.Int.Named, T.Int16.Named, T.Int32.Named, T.Int64.Named, T.Int8.Named:
			code = pkgstrconv.C("FormatInt", T.Int64.From(fAcc), 10)
		case T.Float32.Named:
			code = pkgstrconv.C("FormatFloat", T.Float64.From(fAcc), 'f', -1, 32)
		case T.Float64.Named:
			code = pkgstrconv.C("FormatFloat", fAcc, 'f', -1, 64)
		case T.Rune.Named:
			code, isquoted = pkgstrconv.C("Quote", T.String.From(fAcc)), true
		case T.String.Named:
			code, isquoted = pkgstrconv.C("Quote", fAcc), true
		default:
			code = Call(D(fAcc, N("String")))
		}
	}
	if ensureQuoted && !isquoted {
		code = pkgstrconv.C("Quote", code)
	}
	return
}
