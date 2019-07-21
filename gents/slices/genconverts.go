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

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentConvertMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSlice() {
		if (me.Fields.Add || me.ToMaps.Add) && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of != nil && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName != "" && t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.PkgName == "" {
			if tstruc := ctx.Pkg.Types.Named(t.Expr.GenRef.ArrOrSlice.Of.Pointer.Of.Named.TypeName); tstruc != nil && tstruc.Expr.GenRef.Struct != nil {
				for _, field := range me.Fields.Named {
					if fld := tstruc.Expr.GenRef.Struct.Field(field, false); fld != nil {
						if me.Fields.Add {
							yield.Add(me.genFieldsMethod(t, fld))
						}
						if me.ToMaps.Add {
							yield.Add(me.genToMapMethod(t, fld))
						}
					}
				}
			}
		}
	}
	return
}

func (me *GentConvertMethods) genFieldsMethod(t *gent.Type, field *SynStructField) *SynFunc {
	methodname, tsl :=
		me.Fields.NameWith("field", field.Name), TSlice(field.Type)
	return t.G.T.Method(methodname).Rets(ˇ.R.OfType(tsl)).
		Doc(me.Fields.DocComment.With("N", methodname, "field", field.Name, "T", t.Expr.GenRef.UltimateElemType().String(), "this", Self.Name)).
		Code(
			ˇ.R.Set(B.Make.Of(tsl, B.Len.Of(Self))),
			ForEach(ˇ.I, None, Self,
				ˇ.R.At(ˇ.I).Set(Self.At(ˇ.I).D(N(field.Name))),
			),
		)
}

func (me *GentConvertMethods) genToMapMethod(t *gent.Type, field *SynStructField) *SynFunc {
	methodname, tmap :=
		me.ToMaps.NameWith("field", field.Name), TMap(field.Type, t.Expr.GenRef.ArrOrSlice.Of)
	return t.G.T.Method(methodname).Rets(ˇ.R.OfType(tmap)).
		Doc(me.ToMaps.DocComment.With("N", methodname, "field", field.Name, "T", t.Expr.GenRef.UltimateElemType().String(), "this", Self.Name)).
		Code(
			ˇ.R.Set(B.Make.Of(tmap, B.Len.Of(Self))),
			ForEach(ˇ.I, None, Self,
				ˇ.R.At(Self.At(ˇ.I).D(N(field.Name))).Set(Self.At(ˇ.I)),
			),
		)
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (me *GentConvertMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	me.Fields.Add, me.ToMaps.Add = enabled, enabled
}
