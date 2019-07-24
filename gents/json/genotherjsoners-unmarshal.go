package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func (me *GentTypeJsonMethods) genUnmarshalMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var code ISyn
	self := func() (ISyn, *TypeRef) { return Self, t.Expr.GenRef }
	if t.Expr.GenRef.Struct != nil {
		code = me.genUnmarshalStruct(ctx, self)
	} else if t.Expr.GenRef.ArrOrSlice.Of != nil {
		code = me.genUnmarshalSlice(ctx, self)
	} else if t.Expr.GenRef.Map.OfKey != nil {
		code = me.genUnmarshalMap(ctx, self)
	} else {
		panic(t.Name)
	}
	return t.G.Tª.Method(me.Unmarshal.Name, ˇ.V.OfType(T.SliceOf.Bytes)).Rets(ˇ.Err).
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
	var argtype *TypeRef
	if t.Expr.GenRef.Struct != nil {
		argtype = TMap(T.String, T.Empty.Interface)
		code = me.genUnmarshalStructFromAny(ctx, ˇ.V, Self, t.Expr.GenRef.Struct)
	} else if t.Expr.GenRef.ArrOrSlice.Of != nil {
		argtype = TSlice(T.Empty.Interface)
		code = Block(
			ˇ.Sl.Let(Self.Deref()),
			me.genUnmarshalSliceFromAny(ctx, ˇ.V, ˇ.Sl, t.G.T, t.Expr.GenRef.ArrOrSlice.Of),
			Self.Deref().Set(ˇ.Sl),
		)
	} else if t.Expr.GenRef.Map.OfKey != nil {
		argtype = TMap(T.String, T.Empty.Interface)
		code = Block(
			Var(ˇ.KVs.Name, t.G.T, nil),
			me.genUnmarshalMapFromAny(ctx, ˇ.V, ˇ.KVs, t.G.T, t.Expr.GenRef),
			Self.Deref().Set(ˇ.KVs),
		)
	} else {
		panic(t.Name)
	}
	return t.G.Tª.Method(me.unmarshalHelperMethodName(ctx), ˇ.V.OfType(argtype)).
		Code(code)
}

func (me *GentTypeJsonMethods) genUnmarshalStructFromAny(ctx *gent.Ctx, from ISyn, acc ISyn, t *TypeStruct) (code Syns) {
	for i := range t.Fields {
		fld := &t.Fields[i]
		if jsonname := fld.JsonNameFinal(); jsonname != "" {
			fdacc, nv, nok, nz, ftr :=
				D(acc, N(fld.EffectiveName())), ctx.N("v"), ctx.N("o"), ctx.N("z"), me.typeUnderlyingIfNamed(ctx, fld.Type)
			var tref *TypeRef
			if ftr != nil {
				tref = ftr.Expr.GenRef
			}
			code.Add(
				Tup(nv, nok).Let(At(from, L(jsonname))),
				If(nok.Not(), Then(
					me.genSetToZero(fdacc, nz, fld.Type, tref),
				), me.genUnmarshalSetFromJsonValue(ctx, nv, func() (ISyn, *TypeRef) { return fdacc, fld.Type }),
				),
			)
		}
	}

	return
}

func (*GentTypeJsonMethods) typeUnderlyingIfNamed(ctx *gent.Ctx, t *TypeRef) (tRef *gent.Type) {
	var tnr *TypeRef
	if t.Pointer.Of != nil && t.Pointer.Of.Named.TypeName != "" && t.Pointer.Of.Named.PkgName == "" {
		tnr = t.Pointer.Of
	} else if t.Named.TypeName != "" && t.Named.PkgName == "" {
		tnr = t
	}
	if tnr != nil {
		if tRef = ctx.Pkg.Types.Named(tnr.Named.TypeName); tRef != nil {
			if tRef.Expr.GenRef == nil {
				tRef = nil
			} else {
				for tRef.Expr.GenRef.Named.TypeName != "" && tRef.Expr.GenRef.Named.PkgName == "" {
					if tnext := ctx.Pkg.Types.Named(tRef.Expr.GenRef.Named.TypeName); tnext == nil || tnext.Expr.GenRef == nil {
						break
					} else {
						tRef = tnext
					}
				}
			}
		}
	}
	return
}

func (*GentTypeJsonMethods) genSetToZero(fAcc ISyn, n Named, t *TypeRef, t2 *TypeRef) ISyn {
	return GEN_BYCASE(USUALLY(Block(
		Var(n.Name, t, nil),
		Set(fAcc, n),
	)), UNLESS{
		t.Equiv(T.String) || (t2 != nil && t2.Equiv(T.String)):                                     Set(fAcc, L("")),
		t.Equiv(T.Bool) || (t2 != nil && t2.Equiv(T.Bool)):                                         Set(fAcc, L(false)),
		t.BitSizeIfBuiltInNumberType() != 0 || (t2 != nil && t2.BitSizeIfBuiltInNumberType() != 0): Set(fAcc, L(0)),
		t.CanNil() || (t2 != nil && t2.CanNil()):                                                   Set(fAcc, B.Nil),
	})
}

func (me *GentTypeJsonMethods) genUnmarshalSetFromJsonValue(ctx *gent.Ctx, from ISyn, field func() (ISyn, *TypeRef)) (code Syns) {
	facc, t := field()
	block, nzero := &code, ctx.N("z")
	block.Add(B.Println.Of(from))
	tn := me.typeUnderlyingIfNamed(ctx, t)
	var tnr *TypeRef
	if tn != nil {
		tnr = tn.Expr.GenRef
	}

	setnil := me.genSetToZero(facc, nzero, t, tnr)
	chknil := If(B.Nil.Neq(from), Then(), setnil)
	block.Add(chknil)
	block = &chknil.IfThens[0].Body

	var mslice bool
	var mmap bool
	if tn != nil {
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
	}

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
				C(D(facc, N(me.unmarshalHelperMethodName(ctx))), nv),
			)),
		)

	} else if t.IsBuiltinPrimType(false) || (tnr != nil && tnr.IsBuiltinPrimType(false)) || (tn != nil && tn.IsEnumish()) {
		if t.Equiv(T.Bool) || (tnr != nil && tnr.Equiv(T.Bool)) {
			block.Add(Set(facc, t.From(D(from, T.Bool))))
		} else if t.Equiv(T.String) || (tnr != nil && tnr.Equiv(T.String)) {
			block.Add(Set(facc, t.From(D(from, T.String))))
		} else {
			block.Add(Set(facc, t.From(D(from, T.Float64))))
		}

	} else if tstr := t.Struct; tstr != nil || (tnr != nil && tnr.Struct != nil) {
		if tstr == nil {
			tstr = tnr.Struct
		}
		tmp := ctx.N("t")
		block.Add(
			tmp.Let(D(from, TMap(T.String, T.Empty.Interface))),
			If(B.Nil.Eq(tmp), Then(setnil), me.genUnmarshalStructFromAny(ctx, tmp, facc, tstr)),
		)

	} else if tm := t.Map.OfKey; tm != nil || (tnr != nil && tnr.Map.OfKey != nil) {
		if tm != nil {
			tm = t
		} else {
			tm = tnr
		}
		m := ctx.N("m")
		block.Add(
			m.Let(D(from, TMap(T.String, T.Empty.Interface))),
			If(m.Eq(B.Nil), Then(setnil), me.genUnmarshalMapFromAny(ctx, m, facc, t, tm)),
		)

	} else if tsl := t.ArrOrSlice.Of; tsl != nil || (tnr != nil && tnr.ArrOrSlice.Of != nil) {
		if tsl == nil {
			tsl = tnr.ArrOrSlice.Of
		}
		sl := ctx.N("s")
		block.Add(
			sl.Let(D(from, TSlice(T.Empty.Interface))),
			If(sl.Eq(B.Nil), Then(setnil), me.genUnmarshalSliceFromAny(ctx, sl, facc, t, tsl)),
		)

	}
	return
}

func (me *GentTypeJsonMethods) genUnmarshalMapFromAny(ctx *gent.Ctx, from ISyn, fAcc ISyn, t *TypeRef, tm *TypeRef) (code Syns) {
	mk, mv, tmp := ctx.N("mk"), ctx.N("mv"), ctx.N("t")
	code.Add(
		Set(fAcc, B.Make.Of(t, B.Len.Of(from))),
		ForEach(mk, mv, from,
			Var(tmp.Name, tm.Map.ToVal, nil),
			me.genUnmarshalSetFromJsonValue(ctx, mv,
				func() (ISyn, *TypeRef) { return tmp, tm.Map.ToVal },
			),
			At(fAcc, mk).Set(tmp),
		),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalSliceFromAny(ctx *gent.Ctx, from ISyn, fAcc ISyn, t *TypeRef, tElem *TypeRef) (code Syns) {
	sli, slv := ctx.N("si"), ctx.N("sv")
	code.Add(
		If(B.Len.Of(fAcc).Geq(B.Len.Of(from)), Then(
			Set(fAcc, At(fAcc, Sl(L(0), B.Len.Of(from)))),
		), Else(
			Set(fAcc, B.Make.Of(t, B.Len.Of(from))),
		)),
		ForEach(sli, slv, from,
			me.genUnmarshalSetFromJsonValue(ctx, slv,
				func() (ISyn, *TypeRef) { return At(fAcc, sli), tElem },
			)...,
		),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalMap(ctx *gent.Ctx, _ func() (ISyn, *TypeRef)) (code Syns) {
	code.Add(
		Var(ˇ.KVs.Name, TMap(T.String, T.Empty.Interface), nil),
		ˇ.Err.Set(me.pkgjson.C("Unmarshal", ˇ.V, ˇ.KVs.Addr())),
		If(ˇ.Err.Eq(B.Nil), Then(
			Self.C(me.unmarshalHelperMethodName(ctx), ˇ.KVs),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalSlice(ctx *gent.Ctx, _ func() (ISyn, *TypeRef)) (code Syns) {
	code.Add(
		Var(ˇ.Sl.Name, TSlice(T.Empty.Interface), nil),
		ˇ.Err.Set(me.pkgjson.C("Unmarshal", ˇ.V, ˇ.Sl.Addr())),
		If(ˇ.Err.Eq(B.Nil), Then(
			Self.C(me.unmarshalHelperMethodName(ctx), ˇ.Sl),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalStruct(ctx *gent.Ctx, field func() (ISyn, *TypeRef)) (code Syns) {
	_, t := field()
	code.Add(
		Var(ˇ.KVs.Name, nil, B.Make.Of(TMap(T.String, T.Empty.Interface), len(t.Struct.Fields))),
		ˇ.Err.Set(me.pkgjson.C("Unmarshal", ˇ.V, ˇ.KVs.Addr())),
		If(ˇ.Err.Eq(B.Nil), Then(
			Self.C(me.unmarshalHelperMethodName(ctx), ˇ.KVs),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) unmarshalHelperMethodName(ctx *gent.Ctx) string {
	return ctx.Opt.HelpersPrefix + me.Unmarshal.HelpersPrefix + "FromAny"
}

func (me *GentTypeJsonMethods) unmarshalDecodeMethodName(ctx *gent.Ctx) string {
	if me.Unmarshal.InternalDecodeMethodName != "" {
		return me.Unmarshal.InternalDecodeMethodName
	}
	return ctx.Opt.HelpersPrefix + me.Unmarshal.HelpersPrefix + "Decode"
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeMethod(ctx *gent.Ctx, t *gent.Type, genPanicImpl bool) *SynFunc {
	var code ISyn
	if t.Expr.GenRef.Struct != nil {
		code = me.genUnmarshalStructDecode(ctx, Self, t.Expr.GenRef.Struct)
	} else if t.Expr.GenRef.ArrOrSlice.Of != nil {
		code = Block(ˇ.Err.Set(B.Nil))
		// code = Block(
		// 	ˇ.Sl.Let(Self.Deref()),
		// 	me.genUnmarshalSliceFromAny(ctx, ˇ.V, ˇ.Sl, t.G.T, t.Expr.GenRef.ArrOrSlice.Of),
		// 	Self.Deref().Set(ˇ.Sl),
		// )
	} else if t.Expr.GenRef.Map.OfKey != nil {
		code = Block(
			Var(ˇ.KVs.Name, t.G.T, nil),
			me.genUnmarshalMapDecode(ctx, ˇ.KVs, t.G.T, t.Expr.GenRef),
			If(ˇ.Err.Eq(B.Nil), Then(Self.Deref().Set(ˇ.KVs))),
		)
	} else {
		panic(t.Name)
	}
	return t.G.Tª.Method(me.unmarshalDecodeMethodName(ctx),
		ˇ.J.OfType(TPointer(TFrom(me.pkgjson, "Decoder")))).
		Rets(ˇ.Err).
		Code(code)
}

func (me *GentTypeJsonMethods) genUnmarshalStructDecode(ctx *gent.Ctx, acc ISyn, t *TypeStruct) (code Syns) {
	jk, fn := ctx.N("jk"), ctx.N("fn")
	sw := Switch(fn)
	for i := range t.Fields {
		fld := &t.Fields[i]
		if jsonname := fld.JsonNameFinal(); jsonname != "" {
			sw.Cases.Add(L(jsonname))
		}
	}
	code.Add(me.genUnmarshalDecodeObjOrArr(ctx, '{', jk, fn, Block().Body, Block(sw).Body, Block().Body))

	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeObjOrArr(ctx *gent.Ctx, delim rune, jk Named, k Named, onBeforeLoop Syns, onKey Syns, onSuccess Syns) (code Syns) {
	nexttok, ttok, t2, td :=
		ˇ.J.C("Token"), me.pkgjson.T("Token"), ctx.N("t"), ctx.N("d")
	code.Add(
		Var(t2.Name, ttok, nil),
		Tup(t2, ˇ.Err).Set(nexttok),
		If(ˇ.Err.Eq(B.Nil), Then(
			Tup(td, Nope).Let(t2.D(TFrom(me.pkgjson, "Delim"))),
			If(L(delim).Neq(td), Then(
				ˇ.Err.Set(me.pkgerrs.C("New", "expected "+string(delim))),
			), append(onBeforeLoop,
				For(nil, ˇ.Err.Eq(B.Nil).And(ˇ.J.C("More")), nil,
					Var(jk.Name, ttok, nil),
					Tup(jk, ˇ.Err).Set(nexttok),
					If(ˇ.Err.Eq(B.Nil), append(Then(
						k.Let(jk.D(T.String)),
					), onKey...)),
				),
				If(ˇ.Err.Eq(B.Nil), Then(
					Tup(Nope, ˇ.Err).Set(nexttok),
				)),
				If(ˇ.Err.Eq(B.Nil), onSuccess),
			)),
		)),
	)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalMapDecode(ctx *gent.Ctx, fAcc ISyn, t *TypeRef, tm *TypeRef) (code Syns) {
	nexttok, ttok, jk, mk, jv, t1 :=
		ˇ.J.C("Token"), me.pkgjson.T("Token"), ctx.N("jk"), ctx.N("mk"), ctx.N("jv"), ctx.N("t")
	code.Add(me.genUnmarshalDecodeObjOrArr(ctx, '{', jk, mk, Block(
		t1.Let(B.Make.Of(t)),
	).Body, Block(
		Var(jv.Name, ttok, nil),
		Tup(jv, ˇ.Err).Set(nexttok),
		If(ˇ.Err.Eq(B.Nil), Then(
			B.Println.Of(jv, mk, jk),
		)),
	).Body, Block(
		Set(fAcc, t1),
	).Body)...)
	return
}

func (me *GentTypeJsonMethods) genUnmarshalDecodeAndSet(ctx *gent.Ctx, from ISyn, field func() (ISyn, *TypeRef)) (code Syns) {
	facc, t := field()
	block, nzero := &code, ctx.N("z")
	block.Add(B.Println.Of(from))
	tn := me.typeUnderlyingIfNamed(ctx, t)
	var tnr *TypeRef
	if tn != nil {
		tnr = tn.Expr.GenRef
	}

	setnil := me.genSetToZero(facc, nzero, t, tnr)
	chknil := If(B.Nil.Eq(from), Then(setnil), Else())
	block.Add(chknil)
	block = &chknil.Else.Body

	var mslice bool
	var mmap bool
	if tn != nil {
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
	}

	if mslice || mmap {
		// targ, nv := TSlice(T.Empty.Interface), ctx.N("v")
		// if mmap {
		// 	targ = TMap(T.String, T.Empty.Interface)
		// }
		// ctor := Then()
		// if t.Pointer.Of != nil || (tnr != nil && tnr.Pointer.Of != nil) {
		// 	ctor.Add(Set(facc, B.New.Of(t.Pointer.Of)))
		// } else if t.Map.OfKey != nil || (tnr != nil && tnr.Map.OfKey != nil) {
		// 	ctor.Add(Set(facc, B.Make.Of(t, B.Len.Of(nv))))
		// }
		// block.Add(
		// 	nv.Let(D(from, targ)),
		// 	If(nv.Eq(B.Nil), setnil, Else(
		// 		If(L(len(ctor) > 0).And(B.Nil.Eq(facc)), ctor),
		// 		C(D(facc, N(me.unmarshalDecodeMethodName(ctx))), ˇ.J),
		// 	)),
		// )

	} else if t.IsBuiltinPrimType(false) || (tnr != nil && tnr.IsBuiltinPrimType(false)) || (tn != nil && tn.IsEnumish()) {
		// if t.Equiv(T.Bool) || (tnr != nil && tnr.Equiv(T.Bool)) {
		// 	block.Add(Set(facc, t.From(D(from, T.Bool))))
		// } else if t.Equiv(T.String) || (tnr != nil && tnr.Equiv(T.String)) {
		// 	block.Add(Set(facc, t.From(D(from, T.String))))
		// } else {
		// 	block.Add(Set(facc, t.From(D(from, T.Float64))))
		// }

	} else if tstr := t.Struct; tstr != nil || (tnr != nil && tnr.Struct != nil) {
		// if tstr == nil {
		// 	tstr = tnr.Struct
		// }
		// tmp := ctx.N("t")
		// block.Add(
		// 	tmp.Let(D(from, TMap(T.String, T.Empty.Interface))),
		// 	If(B.Nil.Eq(tmp), Then(setnil), me.genUnmarshalStructFromAny(ctx, tmp, facc, tstr)),
		// )

	} else if tm := t.Map.OfKey; tm != nil || (tnr != nil && tnr.Map.OfKey != nil) {
		// if tm != nil {
		// 	tm = t
		// } else {
		// 	tm = tnr
		// }
		// m := ctx.N("m")
		// block.Add(
		// 	m.Let(D(from, TMap(T.String, T.Empty.Interface))),
		// 	If(m.Eq(B.Nil), Then(setnil), me.genUnmarshalMapFromAny(ctx, m, facc, t, tm)),
		// )

	} else if tsl := t.ArrOrSlice.Of; tsl != nil || (tnr != nil && tnr.ArrOrSlice.Of != nil) {
		// if tsl == nil {
		// 	tsl = tnr.ArrOrSlice.Of
		// }
		// sl := ctx.N("s")
		// block.Add(
		// 	sl.Let(D(from, TSlice(T.Empty.Interface))),
		// 	If(sl.Eq(B.Nil), Then(setnil), me.genUnmarshalSliceFromAny(ctx, sl, facc, t, tsl)),
		// )

	}
	return
}
