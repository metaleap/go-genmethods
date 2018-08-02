package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.Filters.NonNils.NameOrSuffix = "NonNils"
	Gents.Filters.SelectWhere.NameOrSuffix = "SelectWhere"
	Gents.Filters.ByFields.NameOrSuffix = "By{field}"
}

type GentFilteringMethods struct {
	gent.Opts

	NonNils     gent.Variant
	SelectWhere gent.Variant
	ByFields    struct {
		gent.Variant
		Named []string
	}
}

func (this *GentFilteringMethods) genNonNilsMethod(t *gent.Type) *SynFunc {
	return t.G.T.Method(this.NonNils.NameOrSuffix).Rets(ˇ.R.OfType(t.G.T)).
		Doc().
		Code(
			ˇ.R.Set(This),
			For(ˇ.I.Let(0), ˇ.I.Lt(B.Len.Of(ˇ.R)), ˇ.I.Incr1(),
				If(ˇ.R.At(ˇ.I).Eq(B.Nil), Then(
					ˇ.R.Set(B.Append.Of(ˇ.R.Sl(nil, ˇ.I), ˇ.R.Sl(ˇ.I.Plus(1), nil)).Spreads()),
					ˇ.I.Decr1(),
				)),
			),
		)
}

func (this *GentFilteringMethods) genSelectWhereMethod(t *gent.Type) *SynFunc {
	tdpred := TdFunc().Arg("", t.Expr.GenRef.ArrOrSlice.Of).Ret("", T.Bool)
	return t.G.T.Method(this.SelectWhere.NameOrSuffix).Args(ˇ.Ok.OfType(tdpred.T())).Rets(ˇ.R.OfType(t.G.T)).
		Doc().
		Code(
			ˇ.R.Set(B.Make.Of(t.G.T, 0, B.Len.Of(This).Div(2))),
			ForEach(ˇ.I, None, This,
				If(ˇ.Ok.Of(This.At(ˇ.I)), Then(
					ˇ.R.Set(B.Append.Of(ˇ.R, This.At(ˇ.I))),
				)),
			),
		)
}

func (this *GentFilteringMethods) genByFieldMethod(t *gent.Type, field *SynStructField) *SynFunc {
	return t.G.T.Method(this.ByFields.NameWith("field", field.Name)).Args(ˇ.V.OfType(field.Type)).Rets(ˇ.R.OfType(t.Expr.GenRef.ArrOrSlice.Of)).
		Doc().
		Code(
			ForEach(ˇ.I, None, This,
				If(This.At(ˇ.I).D(N(field.Name)).Eq(ˇ.V), Then(
					ˇ.R.Set(This.At(ˇ.I)),
					K.Return,
				)),
			),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentFilteringMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSlice() {
		if this.NonNils.Add && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil {
			yield.Add(this.genNonNilsMethod(t))
		}
		if this.SelectWhere.Add {
			yield.Add(this.genSelectWhereMethod(t))
		}
		if this.ByFields.Add && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName != "" && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.PkgName == "" {
			if tstruc := ctx.Pkg.Types.Named(t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName); tstruc != nil && tstruc.Expr.GenRef.Struct != nil {
				for _, field := range this.ByFields.Named {
					if fld := tstruc.Expr.GenRef.Struct.Field(field, false); fld != nil {
						yield.Add(this.genByFieldMethod(t, fld))
					}
				}
			}
		}
	}
	return
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (this *GentFilteringMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	this.SelectWhere.Add, this.ByFields.Add, this.NonNils.Add = enabled, enabled, enabled
}
