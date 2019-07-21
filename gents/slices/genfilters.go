package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultFiltNonNilsDocComment = "{N} returns only the non-`nil` `{T}` objects contained in `{this}`."
	DefaultFiltFuncDocComment    = "{N} returns only the `{T}` objects contained in `{this}` that satisfy the specified `{ok}` predicate."
	DefaultFiltByDocComment      = "{N} returns {what} `{T}` object(s) encountered in `{this}` whose `{member}` member succeeds for the specified value(s)."
)

func init() {
	Gents.Filters.NonNils.Name, Gents.Filters.NonNils.DocComment = "WhereNotNil", DefaultFiltNonNilsDocComment
	Gents.Filters.Func.Name, Gents.Filters.Func.DocComment = "Where", DefaultFiltFuncDocComment
	Gents.Filters.By.Name, Gents.Filters.By.DocComment = "Where{member}", DefaultFiltByDocComment
}

type GentFilteringMethods struct {
	gent.Opts

	NonNils gent.Variant
	Func    gent.Variant
	By      struct {
		gent.Variation
		Fields  []string
		Methods []NamedTyped
	}
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentFilteringMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSlice() {
		if me.NonNils.Add && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil {
			yield.Add(me.genNonNilsMethod(t))
		}
		if me.Func.Add {
			yield.Add(me.genSelectWhereMethod(t))
		}
		if (!me.By.Disabled) && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName != "" && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.PkgName == "" {
			if tstruc := ctx.Pkg.Types.Named(t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName); tstruc != nil && tstruc.Expr.GenRef.Struct != nil {
				for _, field := range me.By.Fields {
					if fld := tstruc.Expr.GenRef.Struct.Field(field, false); fld != nil {
						yield.Add(me.genByFieldMethod(t, fld))
					}
				}
				for i := range me.By.Methods {
					yield.Add(me.genByMethodMethod(t, i))
				}
			}
		}
	}
	return
}

func (me *GentFilteringMethods) genNonNilsMethod(t *gent.Type) *SynFunc {
	return t.G.T.Method(me.NonNils.Name).Rets(ˇ.R.OfType(t.G.T)).
		Doc(me.NonNils.DocComment.With("N", me.NonNils.Name, "this", Self.Name, "T", t.Expr.GenRef.UltimateElemType().String())).
		Code(
			ˇ.R.Set(Self),
			For(ˇ.I.Let(0), ˇ.I.Lt(B.Len.Of(ˇ.R)), ˇ.I.Incr1(),
				If(ˇ.R.At(ˇ.I).Eq(B.Nil), Then(
					ˇ.R.Set(B.Append.Of(ˇ.R.Sl(nil, ˇ.I), ˇ.R.Sl(ˇ.I.Plus(1), nil)).Spreads()),
					ˇ.I.Decr1(),
				)),
			),
		)
}

func (me *GentFilteringMethods) genSelectWhereMethod(t *gent.Type) *SynFunc {
	tdpred := TdFunc().Arg("", t.Expr.GenRef.ArrOrSlice.Of).Ret("", T.Bool)
	return t.G.T.Method(me.Func.Name).Args(ˇ.Ok.OfType(tdpred.T())).Rets(ˇ.R.OfType(t.G.T)).
		Doc(me.Func.DocComment.With("N", me.Func.Name, "this", Self.Name, "ok", ˇ.Ok.Name, "T", t.Expr.GenRef.UltimateElemType().String())).
		Code(
			ˇ.R.Set(B.Make.Of(t.G.T, 0, B.Len.Of(Self).Div(2))),
			ForEach(ˇ.I, None, Self,
				If(ˇ.Ok.Of(Self.At(ˇ.I)), Then(
					ˇ.R.Set(B.Append.Of(ˇ.R, Self.At(ˇ.I))),
				)),
			),
		)
}

func (me *GentFilteringMethods) genByFieldMethod(t *gent.Type, mem *SynStructField) *SynFunc {
	methodname := me.By.NameWith("member", mem.Name)
	return t.G.T.Method(methodname).
		Args(ˇ.V.OfType(mem.Type)).
		Rets(ˇ.R.OfType(t.Expr.GenRef.ArrOrSlice.Of)).
		Doc(me.By.DocComment.With("N", methodname, "what", "the first", "this", Self.Name, "member", mem.Name, "T", t.Expr.GenRef.UltimateElemType().String())).
		Code(
			ForEach(ˇ.I, None, Self,
				If(Self.At(ˇ.I).D(N(mem.Name)).Eq(ˇ.V), Then(
					ˇ.R.Set(Self.At(ˇ.I)),
					K.Return,
				)),
			),
		)
}

func (me *GentFilteringMethods) genByMethodMethod(t *gent.Type, i int) *SynFunc {
	mem := me.By.Methods[i]
	methodname := me.By.NameWith("member", mem.Name)
	return t.G.T.Method(methodname).
		Args(mem.Type.Func.Args.IfUntypedUse(t.Expr.GenRef.ArrOrSlice.Of)...).
		Rets(ˇ.R.OfType(t.G.T)).
		Doc(me.By.DocComment.With("N", methodname, "what", "only the", "this", Self.Name, "member", mem.Name, "T", t.Expr.GenRef.UltimateElemType().String())).
		Code(
			ˇ.R.Set(Self.C(me.Func.Name, Func("").Args(ˇ.V.OfType(t.Expr.GenRef.ArrOrSlice.Of)).Ret("", T.Bool).Code(
				Ret(ˇ.V.C(mem.Name, mem.Type.Func.Args.Names(false)...)),
			))),
		)
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (me *GentFilteringMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	me.By.Disabled = !enabled
	me.Func.Add, me.NonNils.Add = enabled, enabled
}
