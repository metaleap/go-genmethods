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
	method = t.G.ThisVal.Method(this.MethodName).Sig(&Sigs.NoneToBool).
		Code(
			V.R.SetTo(And(check1, check2)),
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
		invalid1, invalid2, info1, info2, name1, name2 := // 1 refers to enum's first enumerant here, and 2 to last enumerant
			this.IsFirstInvalid, this.IsLastInvalid, "inclusive", "inclusive", t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1]
		if name1 == "_" {
			invalid1, name1 = false, t.Enumish.ConstNames[1]
		}

		var op1, op2 ISyn = V.This.Geq(N(name1)), V.This.Leq(N(name2))
		if invalid1 {
			info1, op1 = "exclusive", V.This.Gt(N(name1))
		}
		if invalid2 {
			info2, op2 = "exclusive", V.This.Lt(N(name2))
		}
		decls = Syns{this.genIsValidMethod(t, op1, op2, name1, info1, name2, info2)}
	}
	return
}
