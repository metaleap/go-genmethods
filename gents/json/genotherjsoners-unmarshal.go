package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func (me *GentTypeJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var code ISyn
	var selftype *TypeRef
	if t.Expr.GenRef.Struct != nil {
		selftype = t.G.Tª
		code = me.genUnmarshalStruct(ctx, func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef })
	} else {
		selftype = t.G.T
		code = Block(ˇ.Err.Set(B.Nil))
	}
	return selftype.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
		Doc(me.Unmarshal.DocComment.With("N", me.Unmarshal.Name)).
		Code(
			GEN_IF(genPanicImpl, Then(
				B.Panic.Of(t.Name),
			), Else(
				code,
			)),
		)
}

func (me *GentTypeJsonMethods) genUnmarshalFromAnyMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var code ISyn
	var selftype, argtype *TypeRef
	if t.Expr.GenRef.Struct != nil {
		selftype, argtype = t.G.Tª, TMap(T.String, T.Empty.Interface)
		code = me.genUnmarshalStructFromAny(ctx, ˇ.V, func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef })
	} else if t.Expr.GenRef.ArrOrSlice.Of != nil {
		selftype, argtype = t.G.T, TSlice(T.Empty.Interface)
		code = Block()
	} else if t.Expr.GenRef.Map.OfKey != nil {
		selftype, argtype = t.G.T, TMap(T.String, T.Empty.Interface)
		code = Block()
	} else {
		panic(t.Name)
	}
	return selftype.Method(me.Unmarshal.HelpersPrefix+"FromAny", ˇ.V.OfType(argtype)).
		Code(code)
}

func (me *GentTypeJsonMethods) genUnmarshalStructFromAny(ctx *gent.Ctx, from ISyn, field func() (ISyn, *TypeRef)) (code Syns) {
	facc, t := field()
	for i := range t.Struct.Fields {
		fld := &t.Struct.Fields[i]
		if jsonname := fld.JsonNameFinal(); jsonname != "" {
			fdacc, nv, nok := D(facc, N(fld.EffectiveName())), ctx.N("v"), ctx.N("o")
			code.Add(
				Tup(nv, nok).Let(At(from, L(jsonname))),
				If(nok,
					me.genUnmarshalSetFromJsonValue(ctx, nv, func() (ISyn, *TypeRef) { return fdacc, fld.Type }),
				),
			)
		}
	}

	return
}

func (me *GentTypeJsonMethods) genUnmarshalSetFromJsonValue(ctx *gent.Ctx, from ISyn, field func() (ISyn, *TypeRef)) (code Syns) {
	facc, t := field()
	block, nzero := &code, ctx.N("z")
	block.Add(B.Println.Of(from))
	var tn *gent.Type
	var tnr *TypeRef
	if t.Pointer.Of != nil && t.Pointer.Of.Named.TypeName != "" && t.Pointer.Of.Named.PkgName == "" {
		tnr = t.Pointer.Of
	} else if t.Named.TypeName != "" && t.Named.PkgName == "" {
		tnr = t
	}
	if tnr != nil {
		if tn = ctx.Pkg.Types.Named(tnr.Named.TypeName); tn != nil {
			if tn.Expr.GenRef == nil {
				tn = nil
			} else {
				for tn.Expr.GenRef.Named.TypeName != "" && tn.Expr.GenRef.Named.PkgName == "" {
					if tnext := ctx.Pkg.Types.Named(tn.Expr.GenRef.Named.TypeName); tnext == nil || tnext.Expr.GenRef == nil {
						break
					} else {
						tn = tnext
					}
				}
			}
		}
	}
	if tnr = nil; tn != nil {
		tnr = tn.Expr.GenRef
	}

	setnil := GEN_BYCASE(USUALLY(Block(
		Var(nzero.Name, t, nil),
		Set(facc, nzero),
	)), UNLESS{
		t.Equiv(T.String) || (tnr != nil && tnr.Equiv(T.String)):                                     Set(facc, L("")),
		t.Equiv(T.Bool) || (tnr != nil && tnr.Equiv(T.Bool)):                                         Set(facc, L(false)),
		t.BitSizeIfBuiltInNumberType() != 0 || (tnr != nil && tnr.BitSizeIfBuiltInNumberType() != 0): Set(facc, L(0)),
		t.CanNil() || (tnr != nil && tnr.CanNil()):                                                   Set(facc, B.Nil),
	})
	chknil := If(B.Nil.Neq(from), Then(), setnil)
	block.Add(chknil)
	block = &chknil.IfThens[0].Body

	if tn != nil {
		var mslice bool
		var mmap bool
		_ = ctx.GentExistsFor(tn, func(g gent.IGent) (ok bool) {
			gjt, ok2 := g.(*GentTypeJsonMethods)
			if ok = ok2 && !gjt.Disabled; ok {
				if mslice = tn.IsSliceOrArray(); !mslice {
					if mmap = tnr.Pointer.Of != nil || tnr.Struct != nil || tnr.Map.OfKey != nil; !mmap {
						panic(tnr.String())
					}
				}
			}
			return
		})
		if mslice || mmap {
			targ, nv := TSlice(T.Empty.Interface), ctx.N("v")
			if mmap {
				targ = TMap(T.String, T.Empty.Interface)
			}
			ctor := Then()
			if t.Pointer.Of != nil || (tnr != nil && tnr.Pointer.Of != nil) {
				ctor.Add(Set(facc, B.New.Of(t.Pointer.Of)))
			} else if t.Map.OfKey != nil || (tnr != nil && tnr.Map.OfKey != nil) {
				ctor.Add(Set(facc, B.Make.Of(t, B.Len.Of(nv))))
			}
			block.Add(
				nv.Let(D(from, targ)),
				If(nv.Eq(B.Nil), setnil, Else(
					If(L(len(ctor) > 0).And(B.Nil.Eq(facc)), ctor),
					C(D(facc, N("jsonUnmarshal_FromAny")), nv),
				)),
			)
			return
		}
	}

	if t.IsBuiltinPrimType(false) || (tnr != nil && tnr.IsBuiltinPrimType(false)) || (tn != nil && tn.IsEnumish()) {
		if t.Equiv(T.Bool) || (tnr != nil && tnr.Equiv(T.Bool)) {
			block.Add(Set(facc, t.From(D(from, T.Bool))))
		} else if t.Equiv(T.String) || (tnr != nil && tnr.Equiv(T.String)) {
			block.Add(Set(facc, t.From(D(from, T.String))))
		} else {
			block.Add(Set(facc, t.From(D(from, T.Float64))))
		}

	} else if tstr := t.Struct; tstr != nil || (tnr != nil && tnr.Struct != nil) {
		tmp := ctx.N("t")
		block.Add(
			tmp.Let(D(from, TMap(T.String, T.Empty.Interface))),
			If(B.Nil.Eq(tmp), setnil, me.genUnmarshalStructFromAny(ctx, tmp, field)),
		)

	} else if tsl := t.ArrOrSlice.Of; tsl != nil || (tnr != nil && tnr.ArrOrSlice.Of != nil) {
		if tsl == nil {
			tsl = tnr.ArrOrSlice.Of
		}
		sl, sli, slv := ctx.N("s"), ctx.N("si"), ctx.N("sv")
		block.Add(
			sl.Let(D(from, TSlice(T.Empty.Interface))),
			If(sl.Eq(B.Nil), Then(
				Set(facc, B.Nil),
			), Else(
				If(B.Len.Of(facc).Geq(B.Len.Of(sl)), Then(
					Set(facc, At(facc, Sl(L(0), B.Len.Of(sl)))),
				), Else(
					Set(facc, B.Make.Of(t, B.Len.Of(sl))),
				)),
				ForEach(sli, slv, sl,
					me.genUnmarshalSetFromJsonValue(ctx, slv,
						func() (ISyn, *TypeRef) { return At(facc, sli), tsl },
					)...,
				),
			)),
		)
	}
	return
}

func (me *GentTypeJsonMethods) genUnmarshalStruct(ctx *gent.Ctx, field func() (ISyn, *TypeRef)) (code Syns) {
	_, t := field()
	ts := t.Struct
	code.Add(
		Var(ˇ.KVs.Name, nil, B.Make.Of(TMap(T.String, T.Empty.Interface), len(ts.Fields))),
		ˇ.Err.Set(ctx.Import("encoding/json").C("Unmarshal", ˇ.V, ˇ.KVs.Addr())),
		If(ˇ.Err.Eq(B.Nil), Then(
			Self.C(me.Unmarshal.HelpersPrefix+"FromAny", ˇ.KVs),
		)),
	)
	return
}
