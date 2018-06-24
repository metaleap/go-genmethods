package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentIndexMethods struct {
	gent.Opts

	IndexOf   IndexMethod
	IndexLast IndexMethod
	IndicesOf struct {
		IndexMethod
		ResultsCapFactor uint
	}
}

type IndexMethod struct {
	Disabled           bool
	DocComment         gent.Str
	Name               string
	VariadicAny        bool
	PredicateVariation struct {
		Disabled bool
		Name     string
	}
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentIndexMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsSliceOrArray() {
		if !this.IndexOf.Disabled {
			decls.Add(this.genIndexOfMethod(t, &this.IndexOf)...)
		}
		if !this.IndexLast.Disabled {
			decls.Add(this.genIndexOfMethod(t, &this.IndexLast)...)
		}
		if !this.IndicesOf.Disabled {
			decls.Add(this.genIndicesMethod(t)...)
		}
	}
	return
}

func (this *GentIndexMethods) genIndexOfMethod(t *gent.Type, self *IndexMethod) (decls Syns) {
	if !self.PredicateVariation.Disabled {

	}
	return
}

func (this *GentIndexMethods) genIndicesMethod(t *gent.Type) (decls Syns) {
	self, ret := &this.IndicesOf, V.R.Typed(T.Sl.Ints)

	gen := func(name string, args NamedsTypeds, predicate ISyn) *SynFunc {
		fn := Fn(t.CodeGen.ThisVal, name, TdFunc(args, ret))
		if self.ResultsCapFactor > 0 {
			fn.Add(Set(V.R, Call(B.Make, ret.Type, L(0), Div(Call(B.Len, V.This), L(self.ResultsCapFactor)))))
		}
		fn.Add(ForRange(V.I, None, V.This,
			If(predicate, Set(V.R, Call(B.Append, V.R, V.I))),
		))
		return fn
	}

	fni := gen(self.Name, NTs("eq", t.Underlying.GenRef.ArrOrSliceOf.Val),
		Eq(I(V.This, V.I), N("eq")))
	decls = append(decls, fni)

	if !self.PredicateVariation.Disabled {
		fnp := gen(self.PredicateVariation.Name, NTs(V.Ok.Name, TrFunc(TdFunc(NTs("", t.Underlying.GenRef.ArrOrSliceOf.Val), NT("", T.Bool)))),
			Call(V.Ok, I(V.This, V.I)))
		decls = append(decls, fnp)
	}
	return
}
