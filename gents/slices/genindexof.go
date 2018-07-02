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
		Variadic bool
	}
	IndexLast struct {
		IndexMethodOpts
		Variadic bool
	}
	IndicesOf struct {
		IndexMethodOpts
		ResultsCapFactor uint
	}
	Contains struct {
		IndexMethodOpts
		VariadicAny bool
		VariadicAll bool
	}
}

type IndexMethodOpts struct {
	Disabled   bool
	DocComment gent.Str
	Name       string
	Predicate  gent.Variant
}

func (this *GentIndexMethods) genIndicesOfMethod(t *gent.Type, methodName string, resultsCapFactor uint, predicate bool) *SynFunc {
	arg, ret := this.indexMethodArg(t, false, predicate), ª.R.T(T.Sl.Ints)
	foreachitemcheckcond := GEN_IF(predicate, Then(
		ª.Ok.Call(ª.This.At(ª.I)), // ok(this[i])
	), Else(
		ª.This.At(ª.I).Eq(ª.V), // this[i] == v
	))

	return t.G.ThisVal.Method(methodName, arg).Rets(ret).
		Doc().
		Code(
			GEN_IF(resultsCapFactor > 0,
				ª.R.Set(C.Make(ret.Type, L(0), C.Len(ª.This).Div(L(resultsCapFactor)))), // r = make([]int, 0, len(this) / ‹resultsCapFactor›)
			),
			ForEach(ª.I, None, ª.This, // for i := range this
				If(foreachitemcheckcond, Then( // if ‹check› {
					ª.R.Set(C.Append(ª.R, ª.I)))), // r = append(r, i)
			),
		)
}

func (this *GentIndexMethods) genIndexOfMethod(t *gent.Type, methodName string, isLast bool, variadic bool, predicate bool) *SynFunc {
	arg := this.indexMethodArg(t, variadic, predicate)
	loopbody := GEN_BYCASE(USUALLY(
		If(ª.This.At(ª.I).Eq(ª.V), Then( // if this[i] == v
			ª.R.Set(ª.I), // r = i
			K.Return)),
	), UNLESS{
		predicate: If(ª.Ok.Call(ª.This.At(ª.I)), Then( // if ok(this[i])
			ª.R.Set(ª.I), // r = i
			K.Return)),
		variadic: ForEach(ª.J, None, arg, // for j := range v
			If(ª.This.At(ª.I).Eq(ª.V.At(ª.J)), Then( // if this[i] == v[j]
				ª.R.Set(ª.I), // r = i
				K.Return))),
	})

	return t.G.ThisVal.Method(methodName, arg).Rets(ª.R.T(T.Int)).
		Doc().
		Code(
			GEN_IF(isLast, Then(
				For(ª.I.Let(C.Len(ª.This).Minus(L(1))), (ª.I.Geq(L(0))), (ª.I.Decr1()), // for i := len(this)-1; i>=0; i--
					loopbody),
			), Else(
				ForEach(ª.I, None, ª.This, // for i := range this
					loopbody),
			)),
			ª.R.Set(-1), // r = -1
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentIndexMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsSliceOrArray() {
		if !this.IndexOf.Disabled {
			decls.Add(this.genIndexOfs(t, false)...)
		}
		if !this.IndexLast.Disabled {
			decls.Add(this.genIndexOfs(t, true)...)
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

func (this *GentIndexMethods) genIndexOfs(t *gent.Type, isLast bool) (decls Syns) {
	self := &this.IndexOf
	if isLast {
		self = &this.IndexLast
	}
	decls.Add(this.genIndexOfMethod(t, self.Name, isLast, self.Variadic, false))
	if self.Predicate.Add {
		decls.Add(this.genIndexOfMethod(t, self.Name+self.Predicate.NameOrSuffix, isLast, false, true))
	}
	return
}

func (this *GentIndexMethods) genIndicesOfs(t *gent.Type) (decls Syns) {
	self := &this.IndicesOf
	decls.Add(this.genIndicesOfMethod(t, self.Name, self.ResultsCapFactor, false))
	if self.Predicate.Add {
		decls.Add(this.genIndicesOfMethod(t, self.Name+self.Predicate.NameOrSuffix, self.ResultsCapFactor, true))
	}
	return
}

func (this *GentIndexMethods) genContainsMethods(t *gent.Type) (decls Syns) {
	return
}

func (this *GentIndexMethods) indexMethodArg(t *gent.Type, variadic bool, predicate bool) (arg NamedTyped) {
	if predicate {
		arg = ª.Ok.T(TdFunc().Arg("", t.Expr.GenRef.ArrOrSliceOf.Val).Ret("", T.Bool).Ref())
	} else if arg = ª.V.T(t.Expr.GenRef.ArrOrSliceOf.Val); variadic {
		arg.Type = TrSlice(arg.Type)
		arg.Type.ArrOrSliceOf.IsEllipsis = true
	}
	return
}
