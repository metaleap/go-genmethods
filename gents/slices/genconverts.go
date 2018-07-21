package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

func init() {
	Gents.Converters.Fields.NameOrSuffix = "All{field}s"
	Gents.Converters.ToMaps.NameOrSuffix = "ToMapBy{field}"
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
	tsl := TrSlice(field.Type)
	return t.G.T.Method(this.Fields.NameWith("field", field.Name)).Rets(ˇ.R.OfType(tsl)).
		Doc().
		Code(
			ˇ.R.Set(B.Make.Of(tsl, B.Len.Of(ˇ.This))),
			ForEach(ˇ.I, None, ˇ.This,
				ˇ.R.At(ˇ.I).Set(ˇ.This.At(ˇ.I).D(N(field.Name))),
			),
		)
}

func (this *GentConvertMethods) genToMapMethod(t *gent.Type, field *SynStructField) *SynFunc {
	tmap := TrMap(field.Type, t.Expr.GenRef.ArrOrSlice.Of)
	return t.G.T.Method(this.ToMaps.NameWith("field", field.Name)).Rets(ˇ.R.OfType(tmap)).
		Doc().
		Code(
			ˇ.R.Set(B.Make.Of(tmap, B.Len.Of(ˇ.This))),
			ForEach(ˇ.I, None, ˇ.This,
				ˇ.R.At(ˇ.This.At(ˇ.I).D(N(field.Name))).Set(ˇ.This.At(ˇ.I)),
			),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentConvertMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSlice() {
		if (this.Fields.Add || this.ToMaps.Add) && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName != "" && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.PkgName == "" {
			if tstruc := ctx.Pkg.Types.Named(t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName); tstruc != nil && tstruc.Expr.GenRef.Struct != nil {
				for _, field := range this.Fields.Named {
					if fld := tstruc.Expr.GenRef.Struct.Field(field); fld != nil {
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
