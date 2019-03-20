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
		DefaultIndexMethodName, DefaultIndicesMethodName, DefaultIndexLastMethodName, DefaultContainsMethodName, true, true, true
	def.IndexOf.Predicate.Name, def.IndicesOf.Predicate.Name, def.IndexLast.Predicate.Name, def.Contains.Predicate.Name =
		DefaultMethodNameSuffixPredicate, DefaultMethodNameSuffixPredicate, DefaultMethodNameSuffixPredicate, DefaultMethodNameSuffixPredicate
}

type GentIndexMethods struct {
	gent.Opts

	IndexOf struct {
		IndexMethodOpts
		Variadic bool
	}

	// `Disabled` in `Gents.IndexOf` by default
	IndexLast struct {
		IndexMethodOpts
		Variadic bool
	}

	// `Disabled` in `Gents.IndexOf` by default
	IndicesOf struct {
		IndexMethodOpts
		ResultsCapFactor uint
	}

	// `Disabled` in `Gents.IndexOf` by default
	Contains struct {
		IndexMethodOpts
		VariadicAny bool
		VariadicAll bool
	}
}

type IndexMethodOpts struct {
	gent.Variation
	Predicate gent.Variant
}

func (me *GentIndexMethods) genIndicesOfMethod(t *gent.Type, methodName string, resultsCapFactor uint, predicate bool) *SynFunc {
	arg, ret := me.indexMethodArg(t, false, predicate), ˇ.R.OfType(T.SliceOf.Ints)
	foreachitemcheckcond := GEN_IF(predicate, Then(
		ˇ.Ok.Of(Self.At(ˇ.I)), // ok(this[i])
	), Else(
		Self.At(ˇ.I).Eq(ˇ.V), // this[i] == v
	))

	return t.G.This.Method(methodName, arg).Rets(ret).
		Doc().
		Code(
			GEN_IF(resultsCapFactor > 0,
				ˇ.R.Set(B.Make.Of(ret.Type, L(0), B.Len.Of(Self).Div(L(resultsCapFactor)))), // r = make([]int, 0, len(this) / ‹resultsCapFactor›)
			),
			ForEach(ˇ.I, None, Self, // for i := range this
				If(foreachitemcheckcond, Then( // if ‹check› {
					ˇ.R.Set(B.Append.Of(ˇ.R, ˇ.I)))), // r = append(r, i)
			),
		)
}

func (me *GentIndexMethods) genIndexOfMethod(t *gent.Type, methodName string, isLast bool, variadic bool, predicate bool) *SynFunc {
	arg, forloop := me.indexMethodArg(t, variadic, predicate), ForEach(ˇ.I, None, Self) // for i := range this
	if isLast {
		forloop = For(ˇ.I.Let(B.Len.Of(Self).Minus(L(1))), (ˇ.I.Gt(-1)), (ˇ.I.Decr1())) // for i := len(this)-1; i > -1; i--
	}

	return t.G.This.Method(methodName, arg).Rets(ˇ.R.OfType(T.Int)).
		Doc().
		Code(
			forloop.Code(GEN_BYCASE(USUALLY(
				If(Self.At(ˇ.I).Eq(ˇ.V), Then( // if this[i] == v
					ˇ.R.Set(ˇ.I), // r = i
					K.Return)),
			), UNLESS{
				predicate: If(ˇ.Ok.Of(Self.At(ˇ.I)), Then( // if ok(this[i])
					ˇ.R.Set(ˇ.I), // r = i
					K.Return)),
				variadic: ForEach(ˇ.J, None, arg, // for j := range v
					If(Self.At(ˇ.I).Eq(ˇ.V.At(ˇ.J)), Then( // if this[i] == v[j]
						ˇ.R.Set(ˇ.I), // r = i
						K.Return))),
			})),
			ˇ.R.Set(-1), // r = -1
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentIndexMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if t.IsSliceOrArray() {
		if !me.IndexOf.Disabled {
			yield.Add(me.genIndexOfs(t, false)...)
		}
		if !me.IndexLast.Disabled {
			yield.Add(me.genIndexOfs(t, true)...)
		}
		if !me.IndicesOf.Disabled {
			yield.Add(me.genIndicesOfs(t)...)
		}
		if !me.Contains.Disabled {
			yield.Add(me.genContainsMethods(t)...)
		}
	}
	return
}

func (me *GentIndexMethods) genIndexOfs(t *gent.Type, isLast bool) (decls Syns) {
	self := &me.IndexOf
	if isLast {
		self = &me.IndexLast
	}
	decls.Add(me.genIndexOfMethod(t, self.Name, isLast, self.Variadic, false))
	if self.Predicate.Add {
		decls.Add(me.genIndexOfMethod(t, self.Name+self.Predicate.Name, isLast, false, true))
	}
	return
}

func (me *GentIndexMethods) genIndicesOfs(t *gent.Type) (decls Syns) {
	self := &me.IndicesOf
	decls.Add(me.genIndicesOfMethod(t, self.Name, self.ResultsCapFactor, false))
	if self.Predicate.Add {
		decls.Add(me.genIndicesOfMethod(t, self.Name+self.Predicate.Name, self.ResultsCapFactor, true))
	}
	return
}

func (me *GentIndexMethods) genContainsMethods(t *gent.Type) (decls Syns) {
	return
}

func (me *GentIndexMethods) indexMethodArg(t *gent.Type, variadic bool, predicate bool) (arg NamedTyped) {
	if predicate {
		arg = ˇ.Ok.OfType(TdFunc().Arg("", t.Expr.GenRef.ArrOrSlice.Of).Ret("", T.Bool).T())
	} else if arg = ˇ.V.OfType(t.Expr.GenRef.ArrOrSlice.Of); variadic {
		arg.Type = TSlice(arg.Type)
		arg.Type.ArrOrSlice.IsEllipsis = true
	}
	return
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (me *GentIndexMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	disabled := !enabled
	me.Contains.Disabled, me.IndexLast.Disabled, me.IndexOf.Disabled, me.IndicesOf.Disabled = disabled, disabled, disabled, disabled
	me.Contains.Predicate.Add, me.IndexLast.Predicate.Add, me.IndexOf.Predicate.Add, me.IndicesOf.Predicate.Add = enabled, enabled, enabled, enabled
}
