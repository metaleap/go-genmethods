package gentstructs

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
)

const (
	DefaultDocCommentTrav = "{N} calls `on` {nf}x: once for each field in this `{T}` with its name, its pointer, `true` if name (or embed name) begins in upper-case (else `false`), and `true` if field is an embed (else `false`)."
	DefaultMethodNameTrav = "StructFieldsTraverse"
)

var (
	on = ˇ.On.OfType(TdFunc(
		T.String.N("name"), T.Empty.Interface.N("ptr"), T.Bool.N("isNameUpperCase"), T.Bool.N("isEmbed")).T())
)

func init() {
	Gents.StructFieldsTrav.DocComment, Gents.StructFieldsTrav.MethodName = DefaultDocCommentTrav, DefaultMethodNameTrav
}

type GentStructFieldsTrav struct {
	gent.Opts

	DocComment      gent.Str
	MethodName      string
	MayIncludeField func(*SynStructField) bool
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentStructFieldsTrav) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if tstruct := t.Expr.GenRef.Struct; tstruct != nil {
		yield.Add(me.genTraverseMethod(ctx, t))
	}
	return
}

func (me *GentStructFieldsTrav) genTraverseMethod(ctx *gent.Ctx, t *gent.Type) *SynFunc {
	numfields := len(t.Expr.GenRef.Struct.Fields)
	return t.G.Tª.Method(me.MethodName, on).
		Doc(me.DocComment.With("N", me.MethodName, "T", t.Name, "nf", ustr.Int(numfields))).
		Code(
			GEN_FOR(0, numfields, 1, func(i int) ISyn {
				fld := &t.Expr.GenRef.Struct.Fields[i]
				if me.MayIncludeField == nil || me.MayIncludeField(fld) {
					fldnamefull := fld.EffectiveName()
					return C(on,
						fldnamefull,
						Self.D(fldnamefull).Addr(),
						GEN_EITHER(ustr.BeginsUpper(fldnamefull), B.True, B.False),
						GEN_EITHER(fld.Name == "", B.True, B.False),
					)
				}
				return Syns{}
			})...,
		)
}
