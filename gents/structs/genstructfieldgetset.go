package gentstructs

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultDocCommentGet = ""
	DefaultMethodNameGet = "StructFieldsGet"
	DefaultDocCommentSet = ""
	DefaultMethodNameSet = "StructFieldsSet"
)

func init() {
	Gents.StructFieldsGetSet.Getter.DocComment, Gents.StructFieldsGetSet.Getter.Name = DefaultDocCommentGet, DefaultMethodNameGet
	Gents.StructFieldsGetSet.Setter.DocComment, Gents.StructFieldsGetSet.Setter.Name = DefaultDocCommentSet, DefaultMethodNameSet
}

type GentStructFieldsGetSet struct {
	gent.Opts

	Getter struct {
		gent.Variation
		ReturnsPtrInsteadOfVal bool
	}
	Setter gent.Variation
}

func (me *GentStructFieldsGetSet) genGetMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Getter.Name).
		Args(ˇ.Name.OfType(T.String), ˇ.V.OfType(T.Empty.Interface)).
		Rets(ˇ.R.OfType(T.Empty.Interface), ˇ.Ok.OfType(T.Bool)).
		Doc().
		Code(
			Switch(ˇ.Name).
				DefaultCase(ˇ.R.Set(ˇ.V)).
				CasesFrom(true, GEN_FOR(0, len(t.Expr.GenRef.Struct.Fields), 1, func(i int) ISyn {
					fldname := t.Expr.GenRef.Struct.Fields[i].EffectiveName()
					return Case(L(fldname),
						ˇ.R.Set(AddrIf(me.Getter.ReturnsPtrInsteadOfVal, Self.D(N(fldname)))),
						ˇ.Ok.Set(B.True),
					)
				})...),
		)
}

func (me *GentStructFieldsGetSet) genSetMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	okN, okT := N("okName"), N("okType")
	return t.G.Tª.Method(me.Setter.Name).
		Args(ˇ.Name.OfType(T.String), ˇ.V.OfType(T.Empty.Interface)).
		Rets(okN.OfType(T.Bool), okT.OfType(T.Bool)).
		Doc().
		Code(
			Switch(ˇ.Name).
				CasesFrom(true, GEN_FOR(0, len(t.Expr.GenRef.Struct.Fields), 1, func(i int) ISyn {
					fld := &t.Expr.GenRef.Struct.Fields[i]
					fldname := fld.EffectiveName()
					if pkgimp := t.SrcFileImportPathByName(fld.Type.Named.PkgName); pkgimp != nil {
						fld.Type.Named.PkgName = string(ctx.Import(pkgimp.ImportPath))
					}
					return Case(L(fldname),
						okN.Set(B.True),
						Tup(ˇ.T, ˇ.Ok).Let(ˇ.V.D(fld.Type)),
						If(ˇ.Ok, Then(okT.Set(B.True), Self.D(fldname).Set(ˇ.T))),
					)
				})...),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentStructFieldsGetSet) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if tstruct := t.Expr.GenRef.Struct; tstruct != nil {
		if !me.Getter.Disabled {
			yield.Add(me.genGetMethod(ctx, t))
		}
		if !me.Setter.Disabled {
			yield.Add(me.genSetMethod(ctx, t))
		}
	}
	return
}
