package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentIndexMethods struct {
	gent.Opts

	IndexOf struct {
		IndexMethodOpts
		gent.Variadic
	}
	IndexLast struct {
		IndexMethodOpts
		gent.Variadic
	}
	IndicesOf struct {
		IndexMethodOpts
		ResultsCapFactor uint
	}
	Contains struct {
		gent.Variant
		gent.Variadic
	}
}

type IndexMethodOpts struct {
	Disabled   bool
	DocComment gent.Str
	Name       string
	Predicate  gent.Variant
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentIndexMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsSliceOrArray() {
		if !this.IndexOf.Disabled {
			decls.Add(this.genIndexOfMethod(t, false)...)
		}
		if !this.IndexLast.Disabled {
			decls.Add(this.genIndexOfMethod(t, true)...)
		}
		if !this.IndicesOf.Disabled {
			decls.Add(this.genIndicesMethod(t)...)
		}

		if this.Contains.Add {
			// 	fn := Fn(t.CodeGen.ThisVal, this.Contains.NameOrSuffix,TrFunc(TdFunc(nil,V.R.Typed(T.Bool))))
		}
	}
	return
}

func (this *GentIndexMethods) genIndexOfMethod(t *gent.Type, isLast bool) (decls Syns) {
	self, ret := &this.IndexOf, V.R.Typed(T.Int)
	if isLast {
		self = &this.IndexLast
	}

	gen := func(name string, arg NamedTyped, stmt ISyn) *SynFunc {
		fn := Fn(t.CodeGen.ThisVal, name, TdFunc(NamedsTypeds{arg}, ret))
		var loop *StmtFor
		if !isLast {
			loop = ForRange(V.I, None, V.This, stmt)
		} else {
			loop = ForLoop(Decl(V.I, Sub(Call(B.Len, V.This), L(1))), Geq(V.I, L(0)), Set(V.I, Sub(V.I, L(1))), stmt)
		}
		fn.Add(loop, Set(V.R, L(-1)))
		return fn
	}

	arg := NT("eq", t.Underlying.GenRef.ArrOrSliceOf.Val)
	var stmt ISyn = If(Eq(I(V.This, V.I), N("eq")), Set(V.R, V.I), K.Ret)
	if self.Variadic {
		arg.Type = TrSlice(arg.Type)
		arg.Type.ArrOrSliceOf.IsEllipsis = true
		stmt = ForRange(V.J, None, arg, If(Eq(I(V.This, V.I), I(N("eq"), V.J)), Set(V.R, V.I), K.Ret))
	}
	fni := gen(self.Name, arg, stmt)
	decls = append(decls, fni)

	if self.Predicate.Add {
		fnp := gen(self.Name+self.Predicate.NameOrSuffix, this.predicateArg(t),
			If(Call(V.Ok, I(V.This, V.I)), Set(V.R, V.I), K.Ret))
		decls = append(decls, fnp)
	}
	return
}

func (this *GentIndexMethods) genIndicesMethod(t *gent.Type) (decls Syns) {
	self, ret := &this.IndicesOf, V.R.Typed(T.Sl.Ints)

	gen := func(name string, arg NamedTyped, predicate ISyn) *SynFunc {
		fn := Fn(t.CodeGen.ThisVal, name, TdFunc(NamedsTypeds{arg}, ret))
		if self.ResultsCapFactor > 0 {
			fn.Add(Set(V.R, Call(B.Make, ret.Type, L(0), Div(Call(B.Len, V.This), L(self.ResultsCapFactor)))))
		}
		fn.Add(ForRange(V.I, None, V.This,
			If(predicate, Set(V.R, Call(B.Append, V.R, V.I))),
		))
		return fn
	}

	fni := gen(self.Name, NT("eq", t.Underlying.GenRef.ArrOrSliceOf.Val),
		Eq(I(V.This, V.I), N("eq")))
	decls = append(decls, fni)

	if self.Predicate.Add {
		fnp := gen(self.Name+self.Predicate.NameOrSuffix, this.predicateArg(t),
			Call(V.Ok, I(V.This, V.I)))
		decls = append(decls, fnp)
	}
	return
}

func (*GentIndexMethods) predicateArg(t *gent.Type) NamedTyped {
	return V.Ok.Typed(TrFunc(TdFunc(NTs("", t.Underlying.GenRef.ArrOrSliceOf.Val), NT("", T.Bool))))
}
