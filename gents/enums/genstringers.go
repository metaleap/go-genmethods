package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
)

const (
	DefaultStringers0DocComments              = "{N} implements the Go standard library's `fmt.Stringer` interface."
	DefaultStringers0MethodName               = "String"
	DefaultStringers1DocComments              = "{N} implements the Go standard library's `fmt.GoStringer` interface."
	DefaultStringers1MethodName               = "GoString"
	DefaultStringersParsersDocComments        = "{N} returns the `{T}` represented by `{s}` (as returned by `{T}.{str}`, {caseSensitivity}), or an `error` if none exists."
	DefaultStringersParsersDocCommentsErrless = "{N} is like `{p}` but returns `{fallback}` for unrecognized inputs."
	DefaultStringersParsersFuncName           = "{T}From{str}"
)

func init() {
	Gents.Stringers.DocComments.Parsers = DefaultStringersParsersDocComments
	Gents.Stringers.DocComments.ParsersErrlessVariant = DefaultStringersParsersDocCommentsErrless
	Gents.Stringers.All = make([]StringMethodOpts, 2)

	Gents.Stringers.All[0].DocComment, Gents.Stringers.All[0].Name =
		DefaultStringers0DocComments, DefaultStringers0MethodName
	Gents.Stringers.All[1].DocComment, Gents.Stringers.All[1].Name, Gents.Stringers.All[1].Disabled =
		DefaultStringers1DocComments, DefaultStringers1MethodName, true

	for i := range Gents.Stringers.All {
		Gents.Stringers.All[i].Parser.FuncName, Gents.Stringers.All[i].Parser.Errless =
			DefaultStringersParsersFuncName, gent.Variant{Name: "Or"}
	}
}

// GentStringersMethods generates for enum type-defs the specified
// `string`ifying methods, optionally with corresponding "parsing" funcs.
//
// An instance with illustrative defaults is in `Gents.String`.
type GentStringersMethods struct {
	gent.Opts

	All         []StringMethodOpts
	DocComments struct {
		Parsers               gent.Str
		ParsersErrlessVariant gent.Str
	}
}

type StringMethodOpts struct {
	gent.Variation
	EnumerantRename gent.Rename
	SkipEarlyChecks bool
	Parser          struct {
		Add               bool
		WithIgnoreCaseCmp bool
		Errless           gent.Variant
		FuncName          gent.Str
	}
}

func (this *StringMethodOpts) genStringerMethod(t *gent.Type, pkgstrconv PkgName, names []string, renames []string) *SynFunc {
	var earlycheck IExprBoolish
	if !this.SkipEarlyChecks {
		earlycheck = Self.Lt(N(names[0])).Or(Self.Gt(N(names[len(names)-1]))) // this < ‹minEnumerant› || this > ‹maxEnumerant›
	}

	ebt := t.Enumish.BaseType
	return t.G.This.Method(this.Name).Rets(ˇ.R.OfType(T.String)).
		Doc(
			this.DocComment.With("N", this.Name, "T", t.Name),
		).
		Code(
			If(earlycheck, GoTo("formatNum")), // if ‹earlycheck› { goto formatNum }
			Switch(Self). // switch this
					DefaultCase(GoTo("formatNum")). // default: goto formatNum
					CasesFrom(true, GEN_FOR(0, len(names), 1, func(i int) ISyn {
					return Case(N(names[i]), // case ‹enumerantIdent›:
						ˇ.R.Set(L(renames[i]))) // r = "‹enumerantNameOrRenamed›"
				})...),
			K.Return,
			Label("formatNum",
				GEN_BYCASE(USUALLY(ˇ.R.Set( // r =
					pkgstrconv.C("FormatInt", T.Int64.From(Self), 10)), // pkg__strconv.FormatInt(int64(this), 10)
				), UNLESS{
					ebt == "int": ˇ.R.Set( // r =
						pkgstrconv.C("Itoa", T.Int.From(Self))), // pkg__strconv.Itoa(int(this))
					ustr.In(ebt, "uint", "uint8", "uint16", "uint32", "uint64"): ˇ.R.Set( // r =
						pkgstrconv.C("FormatUint", T.Uint64.From(Self), 10)), // pkg__strconv.FormatUint(uint64(this), 10)
				}),
			),
		)
}

func (this *StringMethodOpts) genParseFunc(t *gent.Type, docComment gent.Str, pkgstrconv PkgName, pkgstrings PkgName, names []string, renames []string) *SynFunc {
	maybecommonprefix, casehint, parsefuncname :=
		"", "and case-sensitively", this.Parser.FuncName.With("T", t.Name, "str", this.Name)
	if this.Parser.WithIgnoreCaseCmp {
		casehint = "but case-insensitively"
	}

	var earlycheck IExprBoolish
	if minlen, maxlen := ustr.ShortestAndLongest(renames...); !this.SkipEarlyChecks {
		earlycheck = B.Len.Of(ˇ.S).Lt(minlen).Or(B.Len.Of(ˇ.S).Gt(maxlen)) // len(s) < ‹minLen› || len(s) > ‹maxLen›
		if maybecommonprefix = ustr.CommonPrefix(renames...); maybecommonprefix != "" {
			if l := len(maybecommonprefix); !this.Parser.WithIgnoreCaseCmp {
				earlycheck = earlycheck.Or(ˇ.S.Sl(0, l).Neq(maybecommonprefix)) // || s[0:5] != "PREF_" // (for example)
			} else {
				earlycheck = earlycheck.Or(Not(pkgstrings.C("EqualFold", ˇ.S.Sl(0, l), maybecommonprefix))) // || !strings.EqualFold(s[0:5], "Pref_") // (for example)
			}
		}
	}

	ebt, tryparseint := TFrom("", t.Enumish.BaseType), func(inttype *TypeRef, parsefuncname string, args ...IAny) Syns {
		return Syns{
			Var(ˇ.V.Name, inttype, nil), // var v ‹inttype›
			Tup(ˇ.V, ˇ.Err).Set(pkgstrconv.C(parsefuncname, append([]IAny{ˇ.S}, args...)...)), // v, err = strconv.‹ParseFunc›(s, ‹args›)
			If(ˇ.Err.Eq(B.Nil), Then( // if err == nil
				Self.Set(t.G.T.From(ˇ.V)), // this = ‹enumType›(v)
			)),
		}
	}

	return Func(parsefuncname).Args(ˇ.S.OfType(T.String)).Rets(t.G.This, ˇ.Err).
		Doc(
			docComment.With("N", parsefuncname, "T", t.Name, "s", ˇ.S.Name, "str", this.Name, "caseSensitivity", casehint),
		).
		Code(
			If(earlycheck, GoTo("tryParseNum")), // if ‹earlycheck› { goto tryParseNum }
			Block(
				ˇ.T.Let(GEN_EITHER(len(maybecommonprefix) == 0, ˇ.S, ˇ.S.Sl(len(maybecommonprefix), -1))), // t := s   ~|OR|~   t := s[‹lenOfCommonPrefix›:]
				Switch(GEN_EITHER(this.Parser.WithIgnoreCaseCmp, nil, ˇ.T)).
					DefaultCase(GoTo("tryParseNum")). // default: goto tryParseNum
					CasesFrom(true, GEN_FOR(0, len(names), 1, func(i int) ISyn {
						enname, enstrlit := N(names[i]), L(renames[i][len(maybecommonprefix):])
						return Case(GEN_EITHER(!this.Parser.WithIgnoreCaseCmp, enstrlit, pkgstrings.C("EqualFold", ˇ.T, enstrlit)), // case "‹enumerant›"   ~|OR|~   case strings.EqualFold(t, "‹enumerant›")
							Self.Set(enname)) // this = ‹enumerant›
					})...),
				K.Return,
			),
			Label("tryParseNum", GEN_BYCASE(USUALLY(tryparseint(
				T.Int, "Atoi"),
			), UNLESS{
				ustr.In(t.Enumish.BaseType, "int8", "int16", "int32", "int64", "rune"): tryparseint(
					T.Int64, "ParseInt", 10, ebt.BitSizeIfBuiltInNumberType()),
				ustr.In(t.Enumish.BaseType, "uint", "uint8", "uint16", "uint32", "uint64", "byte"): tryparseint(
					T.Uint64, "ParseUint", 10, ebt.BitSizeIfBuiltInNumberType()),
			})),
		)
}

func (this *StringMethodOpts) genParseErrlessFunc(t *gent.Type, docComment gent.Str, funcName string, parseFuncName string) *SynFunc {
	maybe, fallback := N("maybe"+t.Name), N("fallback")
	return Func(funcName).Args(ˇ.S.OfType(T.String), fallback.OfType(t.G.T)).Rets(t.G.This).
		Doc(
			docComment.With("N", funcName, "T", t.Name, "p", parseFuncName, "fallback", fallback.Name),
		).
		Code(
			Tup(maybe, ˇ.Err).Let(C(parseFuncName, ˇ.S)), // maybe,err := ‹parseFunc›(s)
			If(ˇ.Err.Eq(B.Nil), // if err == nil
				Then(Self.Set(maybe)),     // this = maybe
				Else(Self.Set(fallback))), // else this = fallback
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStringersMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if len(this.All) > 0 && t.IsEnumish() {
		yield = make(Syns, 0, 3*len(this.All))
		pkgstrconv, pkgstrings, names := ctx.Import("strconv"), ctx.Import("strings"), make([]string, 0, len(t.Enumish.ConstNames))
		for _, enumerant := range t.Enumish.ConstNames {
			if enumerant != "_" {
				names = append(names, enumerant)
			}
		}
		hadrenameslast, renames := true, make([]string, len(names))

		for i := range this.All {
			if self := &this.All[i]; !self.Disabled {
				if self.EnumerantRename != nil {
					for i := range names {
						renames[i] = self.EnumerantRename(ctx, t, names[i])
					}
					hadrenameslast = true
				} else if hadrenameslast {
					copy(renames, names)
					hadrenameslast = false
				}

				if yield.Add(self.genStringerMethod(t, pkgstrconv, names, renames)); self.Parser.Add {
					fnp := self.genParseFunc(t, this.DocComments.Parsers, pkgstrconv, pkgstrings, names, renames)
					if yield.Add(fnp); self.Parser.Errless.Add {
						yield.Add(self.genParseErrlessFunc(t, this.DocComments.ParsersErrlessVariant, fnp.Name+self.Parser.Errless.Name, fnp.Name))
					}
				}
			}
		}
	}
	return
}

// EnableOrDisableAllVariantsAndOptionals implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStringersMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	disabled := !enabled
	for i := range this.All {
		this.All[i].Disabled = disabled
		this.All[i].Parser.Add = enabled
		this.All[i].Parser.Errless.Add = enabled
	}
}
