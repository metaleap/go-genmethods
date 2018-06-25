package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentIsValidMethod generates a `Valid` method for enum type-defs, which checks
// whether the receiver value seems to be within the range of the known enumerants.
//
// An instance with illustrative defaults is in `Defaults.IsValid`.
type GentIsValidMethod struct {
	gent.Opts

	DocComment     gent.Str
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}

func (this *GentIsValidMethod) genIsValidMethod(t *gent.Type, check1 ISyn, check2 ISyn, name1 string, hint1 string, name2 string, hint2 string) (method *SynFunc) {
	method = t.G.ThisVal.Method(this.MethodName).
		Sig(&Sigs.NoneToBool).
		Code(
			Set(V.R, And(check1, check2)),
		).
		Doc(this.DocComment.With(
			"{N}", this.MethodName,
			"{T}", t.Name,
			"{fn}", name1,
			"{fh}", hint1,
			"{ln}", name2,
			"{lh}", hint2,
		))
	return
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// It returns at most one method if `t` is a suitable enum type-def.
func (this *GentIsValidMethod) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsEnumish() {
		name1, name2, hint1, hint2, invalid1, invalid2 :=
			t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1], "inclusive", "inclusive", this.IsFirstInvalid, this.IsLastInvalid
		if name1 == "_" {
			invalid1, name1 = false, t.Enumish.ConstNames[1]
		}
		var op1, op2 ISyn = Geq(V.This, N(name1)), Leq(V.This, N(name2))
		if invalid1 {
			op1, hint1 = Gt(V.This, N(name1)), "exclusive"
		}
		if invalid2 {
			op2, hint2 = Lt(V.This, N(name2)), "exclusive"
		}
		decls = Syns{this.genIsValidMethod(t, op1, op2, name1, hint1, name2, hint2)}
	}
	return
}
