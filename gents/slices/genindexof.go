package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultIndexMethodName           = "Index"
	DefaultIndicesMethodName         = "Indices"
	DefaultIndexLastMethodName       = "LastIndex"
	DefaultContainsMethodName        = "Contains"
	DefaultMethodNameSuffixPredicate = "Func"
)

func init() {
	def := &Gents.IndexOf
	def.IndexOf.Name, def.IndicesOf.Name, def.IndexLast.Name, def.Contains.Name, def.IndicesOf.Disabled, def.IndexLast.Disabled, def.Contains.Disabled =
		DefaultIndexMethodName, DefaultIndicesMethodName, DefaultIndexLastMethodName, DefaultContainsMethodName, false, false, false
	def.IndexOf.Predicate.NameOrSuffix, def.IndicesOf.Predicate.NameOrSuffix, def.IndexLast.Predicate.NameOrSuffix, def.Contains.Predicate.NameOrSuffix =
		DefaultMethodNameSuffixPredicate, DefaultMethodNameSuffixPredicate, DefaultMethodNameSuffixPredicate, DefaultMethodNameSuffixPredicate
	def.IndexOf.Predicate.Add, def.IndicesOf.Predicate.Add, def.IndexLast.Predicate.Add, def.Contains.Predicate.Add = true, true, true, true
}

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
		IndexMethodOpts
		gent.Variadic
	}
}

type IndexMethodOpts struct {
	Disabled   bool
	DocComment gent.Str
	Name       string
	Predicate  gent.Variant
}

func (*GentIndexMethods) genIndicesOfMethod(t *gent.Type, methodName string, arg NamedTyped, ret NamedTyped, resultsCapFactor uint, perItemPredicateExpr ISyn) *SynFunc {
	return t.G.ThisVal.Method(methodName, arg).Rets(ret).
		Doc().
		Code(
			OnlyIf(resultsCapFactor > 0,
				V.R.SetTo(C.Make(ret.Type, L(0), C.Len(V.This).Div(L(resultsCapFactor)))),
			),
			ForRange(V.I, None, V.This,
				IfThen(perItemPredicateExpr,
					Set(V.R, C.Append(V.R, V.I))),
			),
		)
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
			decls.Add(this.genIndicesOfMethods(t)...)
		}
		if !this.Contains.Disabled {
			decls.Add(this.genContainsMethods(t)...)
		}
	}
	return
}

func (this *GentIndexMethods) genIndexOfMethod(t *gent.Type, isLast bool) (decls Syns) {
	self, ret := &this.IndexOf, V.R.T(T.Int)
	if isLast {
		self = &this.IndexLast
	}

	gen := func(name string, arg NamedTyped, stmt ISyn) *SynFunc {
		fn := Fn(t.G.ThisVal, name, TdFn(Args(arg), ret))
		var loop *StmtFor
		if !isLast {
			loop = ForRange(V.I, None, V.This, stmt)
		} else {
			loop = ForLoop(Decl(V.I, Sub(C.Len(V.This), L(1))), Geq(V.I, L(0)), Set(V.I, Sub(V.I, L(1))), stmt)
		}
		fn.Add(loop, Set(V.R, L(-1)))
		return fn
	}

	arg := this.arg(t, self.Variadic)
	var stmt ISyn = IfThen(Eq(I(V.This, V.I), V.V),
		Set(V.R, V.I), K.Ret)
	if self.Variadic {
		stmt = ForRange(V.J, None, arg,
			IfThen(Eq(I(V.This, V.I), I(V.V, V.J)),
				Set(V.R, V.I), K.Ret))
	}
	fni := gen(self.Name, arg, stmt)
	decls = append(decls, fni)

	if self.Predicate.Add {
		fnp := gen(self.Name+self.Predicate.NameOrSuffix, this.argPredicate(t),
			IfThen(Call(V.Ok, V.This.At(V.I)),
				Set(V.R, V.I), K.Ret))
		decls = append(decls, fnp)
	}
	return
}

func (this *GentIndexMethods) genIndicesOfMethods(t *gent.Type) (decls Syns) {
	self, r := &this.IndicesOf, V.R.T(T.Sl.Ints)
	decls.Add(this.genIndicesOfMethod(t, self.Name, V.V.T(t.Underlying.GenRef.ArrOrSliceOf.Val), r, self.ResultsCapFactor,
		V.This.At(V.I).Eq(V.V),
	))
	if self.Predicate.Add {
		decls.Add(this.genIndicesOfMethod(t, self.Name+self.Predicate.NameOrSuffix, this.argPredicate(t), r, self.ResultsCapFactor,
			Call(V.Ok, V.This.At(V.I)),
		))
	}
	return
}

func (this *GentIndexMethods) genContainsMethods(t *gent.Type) (decls Syns) {
	return
}

func (*GentIndexMethods) arg(t *gent.Type, variadic gent.Variadic) NamedTyped {
	arg := V.V.T(t.Underlying.GenRef.ArrOrSliceOf.Val)
	if variadic {
		arg.Type = TrSlice(arg.Type)
		arg.Type.ArrOrSliceOf.IsEllipsis = true
	}
	return arg
}

func (*GentIndexMethods) argPredicate(t *gent.Type) NamedTyped {
	return V.Ok.T(TdFunc().Arg("", t.Underlying.GenRef.ArrOrSliceOf.Val).Ret("", T.Bool).Ref())
}
