package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentIndexMethods struct {
	gent.Opts

	IndexOf   IndexMethod
	IndexLast IndexMethod
	IndicesOf IndexMethod
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
	self := &this.IndicesOf
	if !self.PredicateVariation.Disabled {
		fn := Fn(t.CodeGen.ThisVal, self.Name, TdFunc(NTs("predicate", TrFunc(TdFunc(NTs("", t.Underlying.GenRef.ArrOrSliceOf.Val), NT("", T.Bool)))), V.Ret.Typed(T.Sl.Ints)),
			K.Ret,
		)
		decls = append(decls, fn)
	}
	return
}
