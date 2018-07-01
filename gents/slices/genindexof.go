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
				ª.R.SetTo(C.Make(ret.Type, L(0), C.Len(ª.This).Div(L(resultsCapFactor)))),
			),
			ForRange(ª.I, None, ª.This,
				IfThen(perItemPredicateExpr,
					Set(ª.R, C.Append(ª.R, ª.I))),
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
			decls.Add(this.genIndicesOfs(t)...)
		}
		if !this.Contains.Disabled {
			decls.Add(this.genContainsMethods(t)...)
		}
	}
	return
}

func (this *GentIndexMethods) genIndexOfMethod(t *gent.Type, isLast bool) (decls Syns) {
	self, ret := &this.IndexOf, ª.R.T(T.Int)
	if isLast {
		self = &this.IndexLast
	}

	gen := func(name string, arg NamedTyped, stmt ISyn) *SynFunc {
		fn := Fn(t.G.ThisVal, name, TdFn(Args(arg), ret))
		var loop *StmtFor
		if !isLast {
			loop = ForRange(ª.I, None, ª.This, stmt)
		} else {
			loop = ForLoop(Decl(ª.I, Sub(C.Len(ª.This), L(1))), Geq(ª.I, L(0)), Set(ª.I, Sub(ª.I, L(1))), stmt)
		}
		fn.Add(loop, Set(ª.R, L(-1)))
		return fn
	}

	arg := this.arg(t, self.Variadic)
	var stmt ISyn = IfThen(Eq(I(ª.This, ª.I), ª.V),
		Set(ª.R, ª.I), K.Return)
	if self.Variadic {
		stmt = ForRange(ª.J, None, arg,
			IfThen(Eq(I(ª.This, ª.I), I(ª.V, ª.J)),
				Set(ª.R, ª.I), K.Return))
	}
	fni := gen(self.Name, arg, stmt)
	decls = append(decls, fni)

	if self.Predicate.Add {
		fnp := gen(self.Name+self.Predicate.NameOrSuffix, this.argPredicate(t),
			IfThen(Call(ª.Ok, ª.This.At(ª.I)),
				Set(ª.R, ª.I), K.Return))
		decls = append(decls, fnp)
	}
	return
}

func (this *GentIndexMethods) genIndicesOfs(t *gent.Type) (decls Syns) {
	self, r := &this.IndicesOf, ª.R.T(T.Sl.Ints)
	decls.Add(this.genIndicesOfMethod(t, self.Name, ª.V.T(t.Underlying.GenRef.ArrOrSliceOf.Val), r, self.ResultsCapFactor,
		ª.This.At(ª.I).Eq(ª.V),
	))
	if self.Predicate.Add {
		decls.Add(this.genIndicesOfMethod(t, self.Name+self.Predicate.NameOrSuffix, this.argPredicate(t), r, self.ResultsCapFactor,
			Call(ª.Ok, ª.This.At(ª.I)),
		))
	}
	return
}

func (this *GentIndexMethods) genContainsMethods(t *gent.Type) (decls Syns) {
	return
}

func (*GentIndexMethods) arg(t *gent.Type, variadic gent.Variadic) NamedTyped {
	arg := ª.V.T(t.Underlying.GenRef.ArrOrSliceOf.Val)
	if variadic {
		arg.Type = TrSlice(arg.Type)
		arg.Type.ArrOrSliceOf.IsEllipsis = true
	}
	return arg
}

func (*GentIndexMethods) argPredicate(t *gent.Type) NamedTyped {
	return ª.Ok.T(TdFunc().Arg("", t.Underlying.GenRef.ArrOrSliceOf.Val).Ret("", T.Bool).Ref())
}
