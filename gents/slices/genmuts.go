package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultMutAppendDocComment = "{N} is a convenience (dot-accessor) short-hand for Go's built-in `append` function."
)

func init() {
	Gents.Mutators.Append.Name, Gents.Mutators.Append.DocComment = "Append", DefaultMutAppendDocComment
}

type GentMutatorMethods struct {
	gent.Opts

	Append gent.Variant
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentMutatorMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSlice() {
		if me.Append.Add {
			yield.Add(me.genAppendMethod(t))
		}
	}
	return
}

func (me *GentMutatorMethods) genAppendMethod(t *gent.Type) *SynFunc {
	return t.G.Tª.Method(me.Append.Name).Args(ˇ.V.OfType(t.Expr.GenRef.ArrOrSlice.Of)).Spreads().
		Doc(me.Append.DocComment.With("N", me.Append.Name)).
		Code(
			Self.Deref().Set(B.Append.Of(Self.Deref(), ˇ.V).Spreads()),
		)
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (me *GentMutatorMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	me.Append.Add = enabled
}
