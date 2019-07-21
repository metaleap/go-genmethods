package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultIsValidDocComment = "{N} returns whether the value of this `{T}` is between `{fn}` ({fh}) and `{ln}` ({lh})."
	DefaultIsValidMethodName = "Valid"
)

func init() {
	Gents.IsValid.DocComment, Gents.IsValid.MethodName = DefaultIsValidDocComment, DefaultIsValidMethodName
}

// GentIsValidMethod generates a `Valid` method for enum type-defs, which checks
// whether the receiver value seems to be within the range of the known enumerants.
//
// An instance with illustrative defaults is in `Gents.IsValid`.
type GentIsValidMethod struct {
	gent.Opts

	DocComment     gent.Str
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}

type comparisonOperator = func(IAny) IExprBoolish

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// It returns at most one method if `t` is a suitable enum type-def.
func (me *GentIsValidMethod) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsEnumish() {
		invalid1, invalid2, info1, info2, name1, name2 := // 1 refers to enum's first enumerant here, and 2 to last enumerant
			me.IsFirstInvalid, me.IsLastInvalid, "inclusive", "inclusive", t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1]
		if name1 == "_" {
			invalid1, name1 = false, t.Enumish.ConstNames[1]
		}

		var op1, op2 comparisonOperator = Self.Geq, Self.Leq
		if invalid1 {
			info1, op1 = "exclusive", Self.Gt
		}
		if invalid2 {
			info2, op2 = "exclusive", Self.Lt
		}
		yield = Syns{me.genIsValidMethod(t, op1, op2, name1, name2, info1, info2)}
	}
	return
}

func (me *GentIsValidMethod) genIsValidMethod(t *gent.Type, gtOrGeq, ltOrLeq comparisonOperator, name1, name2 string, hint1, hint2 string) *SynFunc {
	return t.G.This.Method(me.MethodName).Rets(ˇ.R.OfType(T.Bool)).
		Doc(me.DocComment.With(
			"N", me.MethodName, "T", t.Name,
			"fn", name1, "fh", hint1, "ln", name2, "lh", hint2,
		)).
		Code(
			ˇ.R.Set(gtOrGeq(N(name1)).And(ltOrLeq(N(name2)))), // r = (this >? ‹lowestEnumerant›) && (this <? ‹highestEnumerant›)
		)
}
