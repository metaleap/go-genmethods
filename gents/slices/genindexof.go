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
	arg, ret := this.indexMethodArg(t, false, predicate), ˇ.R.OfType(T.Sl.Ints)
	foreachitemcheckcond := GEN_IF(predicate, Then(
		ˇ.Ok.Call(ˇ.This.At(ˇ.I)), // ok(this[i])
	), Else(
		ˇ.This.At(ˇ.I).Eq(ˇ.V), // this[i] == v
	))

	return t.G.ThisVal.Method(methodName, arg).Rets(ret).
		Doc().
		Code(
			GEN_IF(resultsCapFactor > 0,
				ˇ.R.Set(C.Make(ret.Type, L(0), C.Len(ˇ.This).Div(L(resultsCapFactor)))), // r = make([]int, 0, len(this) / ‹resultsCapFactor›)
			),
			ForEach(ˇ.I, None, ˇ.This, // for i := range this
				If(foreachitemcheckcond, Then( // if ‹check› {
					ˇ.R.Set(C.Append(ˇ.R, ˇ.I)))), // r = append(r, i)
			),
		)
}

func (this *GentIndexMethods) genIndexOfMethod(t *gent.Type, methodName string, isLast bool, variadic bool, predicate bool) *SynFunc {
	arg := this.indexMethodArg(t, variadic, predicate)
	loopbody := GEN_BYCASE(USUALLY(
		If(ˇ.This.At(ˇ.I).Eq(ˇ.V), Then( // if this[i] == v
			ˇ.R.Set(ˇ.I), // r = i
			K.Return)),
	), UNLESS{
		predicate: If(ˇ.Ok.Call(ˇ.This.At(ˇ.I)), Then( // if ok(this[i])
			ˇ.R.Set(ˇ.I), // r = i
			K.Return)),
		variadic: ForEach(ˇ.J, None, arg, // for j := range v
			If(ˇ.This.At(ˇ.I).Eq(ˇ.V.At(ˇ.J)), Then( // if this[i] == v[j]
				ˇ.R.Set(ˇ.I), // r = i
				K.Return))),
	})

	return t.G.ThisVal.Method(methodName, arg).Rets(ˇ.R.OfType(T.Int)).
		Doc().
		Code(
			GEN_IF(isLast, Then(
				For(ˇ.I.Let(C.Len(ˇ.This).Minus(L(1))), (ˇ.I.Geq(L(0))), (ˇ.I.Decr1()), // for i := len(this)-1; i>=0; i--
					loopbody),
			), Else(
				ForEach(ˇ.I, None, ˇ.This, // for i := range this
					loopbody),
			)),
			ˇ.R.Set(-1), // r = -1
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
		arg = ˇ.Ok.OfType(TdFunc().Arg("", t.Expr.GenRef.ArrOrSlice.Of).Ret("", T.Bool).Ref())
	} else if arg = ˇ.V.OfType(t.Expr.GenRef.ArrOrSlice.Of); variadic {
		arg.Type = TrSlice(arg.Type)
		arg.Type.ArrOrSlice.IsEllipsis = true
	}
	return
}
