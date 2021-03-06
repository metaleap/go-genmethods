package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
)

func (me *GentTypeJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	return t.G.Tª.Method(me.Unmarshal.Name, ˇ.B.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				B.Panic.Of(t.Name),
			), Else(
				ˇ.J.Let(me.pkgjson.C("NewDecoder", me.pkgbytes.C("NewReader", ˇ.B))),
				ˇ.J.C("UseNumber"),
				ˇ.Err.Set(Self.C(me.unmarshalDecodeMethodName(ctx), ˇ.J)),
			)),
		)
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var self ISyn = Self
	if t.Expr.GenRef.ArrOrSlice.Of != nil || t.Expr.GenRef.Map.OfKey != nil {
		self = Deref(self)
	}
	return t.G.Tª.Method(me.unmarshalDecodeMethodName(ctx),
		ˇ.J.OfType(me.pkgjson.Tª("Decoder"))).
		Rets(ˇ.Err).
		Code(me.genUnmarshalDecodeBasedOnType(ctx, self, t.Expr.GenRef, true)...)
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeBasedOnType(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, canUseExtraDefs bool) Syns {
	if canUseExtraDefs {
		me.genUnmarshalExtraDefs(ctx)
		for _, tref := range me.Unmarshal.CommonTypesToExtractToHelpers {
			if tref.Equiv(fType) {
				return Syns{
					Tup(fAcc, ˇ.Err).Set(C(me.unmarshalExtraDefName(ctx, tref), ˇ.J)),
				}
			}
		}
	}
	switch {
	case fType.ArrOrSlice.Of != nil:
		return me.genUnmarshalDecodeSlice(ctx, fAcc, fType, false)
	case fType.Map.OfKey != nil:
		return me.genUnmarshalDecodeMap(ctx, fAcc, fType, false)
	case fType.Pointer.Of != nil:
		return me.genUnmarshalDecodePtr(ctx, fAcc, fType)
	case fType.Struct != nil:
		return me.genUnmarshalDecodeStruct(ctx, fAcc, fType)
	case fType.IsBuiltinPrimType(false):
		return me.genUnmarshalDecodeBuiltinPrim(ctx, fAcc, fType)
	case fType.Interface != nil:
		return me.genUnmarshalDecodeIface(ctx, fAcc, fType)
	case fType.Named.TypeName != "":
		var pkg *gent.Pkg
		if fType.Named.PkgName == "" {
			pkg = ctx.Pkg
		} else { // if pkg = gent.TryExtPkg(t.Named.PkgName); pkg == nil  /* ext-pkgs stuff not really working just yet, TODO when it becomes more pressing */ {
			return me.genUnmarshalDecodeUnknown(ctx, fAcc, fType)
		}

		if gt := pkg.Types.Named(fType.Named.TypeName); gt == nil {
			panic(fType.Named.TypeName)
		} else {
			if gt.IsEnumish() {
				return me.genUnmarshalDecodeBuiltinPrim(ctx, fAcc, fType)
			} else if ctx.GentExistsFor(gt, func(g gent.IGent) bool {
				gjt, ok := g.(*GentTypeJsonMethods)
				return ok && !(gjt.Disabled || gjt.Unmarshal.Disabled)
			}) {
				return Block(
					ˇ.Err.Set(C(D(fAcc, N(me.unmarshalDecodeMethodName(ctx))), ˇ.J)),
				).Body
			} else {
				return me.genUnmarshalDecodeBasedOnType(ctx, fAcc, gt.Expr.GenRef, true)
			}
		}
	default:
		panic(fType.String())
	}
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeSlice(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, skipDelim bool) (code Syns) {
	sl, val := ctx.N("sl"), ctx.N("v")
	code.Add(me.genUnmarshalDecodeObjOrArr(ctx, '[', skipDelim, None, None, Block(
		sl.Let(B.Make.Of(fType, 0, me.Unmarshal.DefaultCaps.Slices)),
	).Body, Block(
		Var(val.Name, fType.ArrOrSlice.Of, nil),
		me.genUnmarshalDecodeBasedOnType(ctx, val, fType.ArrOrSlice.Of, true),
		If(ˇ.Err.Eq(B.Nil), Then(
			sl.Set(B.Append.Of(sl, val)),
		)),
	).Body, Block(
		Set(fAcc, sl),
	).Body))
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeMap(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef, skipDelim bool) (code Syns) {
	jk, mk, mv, t1 := ctx.N("jk"), ctx.N("mk"), ctx.N("mv"), ctx.N("t")
	code.Add(me.genUnmarshalDecodeObjOrArr(ctx, '{', skipDelim, jk, mk, Block(
		t1.Let(B.Make.Of(fType, me.Unmarshal.DefaultCaps.Maps)),
	).Body, Block(
		Var(mv.Name, fType.Map.ToVal, nil),
		me.genUnmarshalDecodeBasedOnType(ctx, mv, fType.Map.ToVal, true),
		If(ˇ.Err.Eq(B.Nil), Then(
			At(t1, mk).Set(mv),
		)),
	).Body, Block(
		Set(fAcc, t1),
	).Body)...)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeStruct(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (code Syns) {
	jk, fn := ctx.N("jk"), ctx.N("fn")
	fieldnamecases := Switch(fn)
	fieldnamecases.Cases = me.genUnmarshalDecodeStructFieldNameCases(ctx, fAcc, fType)
	code = me.genUnmarshalDecodeObjOrArr(ctx, '{', false, jk, fn, nil, Syns{fieldnamecases}, nil)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeStructFieldNameCases(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (fieldNameCases SynCases) {
	for i := range fType.Struct.Fields {
		fld := &fType.Struct.Fields[i]
		if ft := ctx.Pkg.Types.Named(fld.Type.UltimateElemType().Named.TypeName); fld.Name == "" && ft != nil && ft.Expr.GenRef.Struct != nil {
			fieldNameCases = append(fieldNameCases, me.genUnmarshalDecodeStructFieldNameCases(ctx, D(fAcc, N(fld.EffectiveName())), ft.Expr.GenRef)...)
		} else if jsonfieldname := fld.JsonNameFinal(); jsonfieldname != "" {
			fieldNameCases.Add(L(jsonfieldname),
				me.genUnmarshalDecodeBasedOnType(ctx, D(fAcc, N(fld.EffectiveName())), fld.Type, true)...,
			)
		}
	}
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodePtr(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (code Syns) {
	pv := ctx.N("pv")
	code.Add(
		Var(pv.Name, fType.Pointer.Of, nil),
		me.genUnmarshalDecodeBasedOnType(ctx, pv, fType.Pointer.Of, true),
		If(ˇ.Err.Eq(B.Nil), Then(
			Set(fAcc, pv.Addr()),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeBuiltinPrim(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (code Syns) {
	var valtype, valmethodtype *TypeRef
	var valmethodname string

	switch fType.Named {
	case T.String.Named:
		valtype = fType
	case T.Bool.Named:
		valtype = fType
	case T.Byte.Named, T.Uint.Named, T.Uint16.Named, T.Uint32.Named, T.Uint64.Named, T.Uint8.Named:
		valtype, valmethodname, valmethodtype = me.pkgjson.T("Number"), "Int64", T.Int64
	case T.Float32.Named, T.Float64.Named:
		valtype, valmethodname, valmethodtype = me.pkgjson.T("Number"), "Float64", T.Float64
	case T.Rune.Named, T.Int.Named, T.Int16.Named, T.Int32.Named, T.Int64.Named, T.Int8.Named:
		valtype, valmethodname, valmethodtype = me.pkgjson.T("Number"), "Int64", T.Int64
	default:
		var t *gent.Type
		if fType.Named.TypeName != "" && fType.Named.PkgName == "" {
			t = ctx.Pkg.Types.Named(fType.Named.TypeName)
		}
		if t == nil {
			panic(fType)
		} else {
			valtype, valmethodname, valmethodtype = me.pkgjson.T("Number"), "Int64", T.Int64
		}
	}
	tok, tmp := ctx.N("jt"), ctx.N("v")
	code.Add(
		Var(tok.Name, me.pkgjson.T("Token"), nil),
		Tup(tok, ˇ.Err).Set(ˇ.J.C("Token")),
		If(ˇ.Err.Eq(B.Nil).And(tok.Neq(B.Nil)), Then(
			Switch(ˇ.V.Let(tok.D("type"))).
				Case(B.Nil).
				Case(valtype,
					GEN_IF(valmethodname == "", Then(
						Set(fAcc, ˇ.V),
					), Else(
						Var(tmp.Name, valmethodtype, nil),
						Tup(tmp, ˇ.Err).Set(ˇ.V.C(valmethodname)),
						If(ˇ.Err.Eq(B.Nil), Then(
							Set(fAcc, fType.From(tmp)),
						)),
					)),
				).
				DefaultCase(
					ˇ.Err.Set(me.pkgerrs.C("New", "expected "+valtype.String())),
				),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeUnknown(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (code Syns) {
	tmp := ctx.N("ix")
	code.Add(me.Unmarshal.OnStdlibFallbacks(ctx, fAcc,
		Var(tmp.Name, fType, nil),
		ˇ.Err.Set(ˇ.J.C("Decode", tmp.Addr())),
		If(ˇ.Err.Eq(B.Nil), Then(
			Set(fAcc, tmp),
		)),
	))
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeIface(ctx *gent.Ctx, fAcc ISyn, fType *TypeRef) (code Syns) {
	me.genUnmarshalExtraDefs(ctx)
	tok, m, sl := ctx.N("ttt"), ctx.N("mmm"), ctx.N("slsl")
	code.Add(
		Var(tok.Name, me.pkgjson.T("Token"), nil),
		Tup(tok, ˇ.Err).Set(ˇ.J.C("Token")),
		If(ˇ.Err.Eq(B.Nil).And(tok.Neq(B.Nil)), Then(
			Switch(ˇ.V.Let(tok.D("type"))).
				Case(B.Nil).
				Case(T.String,
					Set(fAcc, ˇ.V)).
				Case(T.Bool,
					Set(fAcc, ˇ.V)).
				Case(me.pkgjson.T("Number"),
					Tup(fAcc, ˇ.Err).Set(ˇ.V.C("Float64"))).
				Case(me.pkgjson.T("Delim"),
					Switch(ˇ.V).
						Case(L('{'),
							Var(m.Name, TMap(T.String, T.Empty.Interface), nil),
							me.genUnmarshalDecodeMap(ctx, m, TMap(T.String, T.Empty.Interface), true),
							If(ˇ.Err.Eq(B.Nil),
								Set(fAcc, m),
							),
						).
						Case(L('['),
							Var(sl.Name, TSlice(T.Empty.Interface), nil),
							me.genUnmarshalDecodeSlice(ctx, sl, TSlice(T.Empty.Interface), true),
							If(ˇ.Err.Eq(B.Nil),
								Set(fAcc, sl),
							),
						),
				),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeObjOrArr(ctx *gent.Ctx, delim byte, skipDelim bool, jk Named, k Named, onBeforeLoop Syns, onNextValue Syns, onSuccess Syns) (code Syns) {
	nexttok, ttok, t2, td :=
		ˇ.J.C("Token"), me.pkgjson.T("Token"), ctx.N("t"), ctx.N("d")
	proceedinto := append(onBeforeLoop,
		For(nil, ˇ.Err.Eq(B.Nil).And(ˇ.J.C("More")), nil,
			GEN_IF(jk.Name == "", onNextValue, Else(
				Var(jk.Name, ttok, nil),
				Tup(jk, ˇ.Err).Set(nexttok),
				If(ˇ.Err.Eq(B.Nil), append(Then(
					k.Let(jk.D(T.String)),
				), onNextValue...)),
			)),
		),
		If(ˇ.Err.Eq(B.Nil), Then(
			Tup(Nope, ˇ.Err).Set(nexttok),
		)),
		If(ˇ.Err.Eq(B.Nil),
			onSuccess,
		),
	)
	if skipDelim {
		code.Add(proceedinto...)
	} else {
		code.Add(
			Var(t2.Name, ttok, nil),
			Tup(t2, ˇ.Err).Set(nexttok),
			If(ˇ.Err.Eq(B.Nil).And(t2.Neq(B.Nil)), Then(
				Switch(td.Let(t2.D(N("type")))).
					Case(B.Nil).
					Case(me.pkgjson.T("Delim"),
						If(L(delim).Neq(td), Then(
							ˇ.Err.Set(me.pkgerrs.C("New", "expected "+string(delim))),
						), proceedinto)).
					DefaultCase(
						ˇ.Err.Set(me.pkgerrs.C("New", "expected "+string(delim))),
					),
			),
			),
		)
	}
	return
}

func (me *GentTypeJsonMethods) unmarshalDecodeMethodName(ctx *gent.Ctx) string {
	if me.Unmarshal.InternalDecodeMethodName != "" {
		return me.Unmarshal.InternalDecodeMethodName
	}
	return ctx.Opt.HelpersPrefix + me.Unmarshal.HelpersPrefix + "Decode"
}

func (me *GentTypeJsonMethods) unmarshalExtraDefName(ctx *gent.Ctx, t *TypeRef) string {
	return ctx.Opt.HelpersPrefix + me.Unmarshal.HelpersPrefix + ustr.ReplB(t.String(), '[', 's', ']', '_', '*', 'p', '{', '_', '}', '_', '.', '_')
}

func (me *GentTypeJsonMethods) genUnmarshalExtraDefs(ctx *gent.Ctx) {
	if !me.Unmarshal.commonTypesToExtraDefsDone {
		me.Unmarshal.commonTypesToExtraDefsDone = true
		defsdone := map[string]struct{}{}
		for _, ftype := range me.Unmarshal.CommonTypesToExtractToHelpers {
			defname := me.unmarshalExtraDefName(ctx, ftype)
			if _, defdone := defsdone[defname]; !defdone {
				defsdone[defname] = struct{}{}
				ctx.ExtraDefs = append(ctx.ExtraDefs, Func(defname, ˇ.J.OfType(me.pkgjson.Tª("Decoder"))).
					Rets(ˇ.R.OfType(ftype), ˇ.Err).
					Code(me.genUnmarshalDecodeBasedOnType(ctx, ˇ.R, ftype, false)...),
				)
			}
		}
	}
}
