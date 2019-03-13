package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultConvFieldsDocComment = "{N} returns all `{field}` values of the constituent `{T}`s in `{this}`."
	DefaultConvToMapsDocComment = "{N} converts `{this}` into a `map` indexed by the `{field}` values of its constituent `{T}`s."
)

func init() {
	Gents.Converters.Fields.Name, Gents.Converters.Fields.DocComment = "All{field}s", DefaultConvFieldsDocComment
	Gents.Converters.ToMaps.Name, Gents.Converters.ToMaps.DocComment = "ToMapBy{field}", DefaultConvToMapsDocComment
}

type GentConvertMethods struct {
	gent.Opts

	Fields struct {
		gent.Variant
		Named []string
	}
	ToMaps struct {
		gent.Variant
		ByFields []string
	}
}

func (this *GentConvertMethods) genFieldsMethod(t *gent.Type, field *SynStructField) *SynFunc {
	methodname, tsl :=
		this.Fields.NameWith("field", field.Name), TSlice(field.Type)
	return t.G.T.Method(methodname).Rets(ˇ.R.OfType(tsl)).
		Doc(this.Fields.DocComment.With("N", methodname, "field", field.Name, "T", t.Expr.GenRef.UltimateElemType().String(), "this", Self.Name)).
		Code(
			ˇ.R.Set(B.Make.Of(tsl, B.Len.Of(Self))),
			ForEach(ˇ.I, None, Self,
				ˇ.R.At(ˇ.I).Set(Self.At(ˇ.I).D(N(field.Name))),
			),
		)
}

func (this *GentConvertMethods) genToMapMethod(t *gent.Type, field *SynStructField) *SynFunc {
	methodname, tmap :=
		this.ToMaps.NameWith("field", field.Name), TMap(field.Type, t.Expr.GenRef.ArrOrSlice.Of)
	return t.G.T.Method(methodname).Rets(ˇ.R.OfType(tmap)).
		Doc(this.ToMaps.DocComment.With("N", methodname, "field", field.Name, "T", t.Expr.GenRef.UltimateElemType().String(), "this", Self.Name)).
		Code(
			ˇ.R.Set(B.Make.Of(tmap, B.Len.Of(Self))),
			ForEach(ˇ.I, None, Self,
				ˇ.R.At(Self.At(ˇ.I).D(N(field.Name))).Set(Self.At(ˇ.I)),
			),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentConvertMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSlice() {
		if (this.Fields.Add || this.ToMaps.Add) && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName != "" && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.PkgName == "" {
			if tstruc := ctx.Pkg.Types.Named(t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName); tstruc != nil && tstruc.Expr.GenRef.Struct != nil {
				for _, field := range this.Fields.Named {
					if fld := tstruc.Expr.GenRef.Struct.Field(field, false); fld != nil {
						if this.Fields.Add {
							yield.Add(this.genFieldsMethod(t, fld))
						}
						if this.ToMaps.Add {
							yield.Add(this.genToMapMethod(t, fld))
						}
					}
				}
			}
		}
	}
	return
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (this *GentConvertMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	this.Fields.Add, this.ToMaps.Add = enabled, enabled
}
