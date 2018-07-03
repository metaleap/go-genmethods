package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
)

const (
	DefaultStringers0DocComments              = "{N} implements the `fmt.Stringer` interface."
	DefaultStringers0MethodName               = "String"
	DefaultStringers1DocComments              = "{N} implements the `fmt.GoStringer` interface."
	DefaultStringers1MethodName               = "GoString"
	DefaultStringersParsersDocComments        = "{N} returns the `{T}` represented by `{s}` (as returned by `{T}.{str}`, {caseSensitivity}), or an `error` if none exists."
	DefaultStringersParsersDocCommentsErrless = "{N} is like `{p}` but returns `{fallback}` for unrecognized inputs."
	DefaultStringersParsersFuncName           = "{T}From{str}"
)

func init() {
	Gents.Stringers.Stringers = []StringMethodOpts{
		{DocComment: DefaultStringers0DocComments, Name: DefaultStringers0MethodName, EnumerantRename: nil},
		{DocComment: DefaultStringers1DocComments, Name: DefaultStringers1MethodName, Disabled: true},
	}
	for i := range Gents.Stringers.Stringers {
		Gents.Stringers.Stringers[i].Parser.Add, Gents.Stringers.Stringers[i].Parser.FuncName, Gents.Stringers.Stringers[i].Parser.Errless =
			true, DefaultStringersParsersFuncName, gent.Variant{Add: false, NameOrSuffix: "Or"}
	}
	Gents.Stringers.DocComments.Parsers = DefaultStringersParsersDocComments
	Gents.Stringers.DocComments.ParsersErrlessVariant = DefaultStringersParsersDocCommentsErrless
}

// GentStringersMethods generates for enum type-defs the specified
// `string`ifying methods, optionally with corresponding "parsing" funcs.
//
// An instance with illustrative defaults is in `Gents.String`.
type GentStringersMethods struct {
	gent.Opts

	Stringers   []StringMethodOpts
	DocComments struct {
		Parsers               gent.Str
		ParsersErrlessVariant gent.Str
	}
}

type StringMethodOpts struct {
	Disabled        bool
	DocComment      gent.Str
	Name            string
	EnumerantRename func(string) string
	Parser          struct {
		Add               bool
		WithIgnoreCaseCmp bool
		Errless           gent.Variant
		FuncName          gent.Str
	}
}

func (this *StringMethodOpts) genStringerMethod(t *gent.Type, pkgstrconv PkgName, namesAndRenames []Any) *SynFunc {
	ebt := t.Enumish.BaseType
	return t.G.This.Method(this.Name).Rets(ˇ.R.OfType(T.String)).
		Doc(
			this.DocComment.With("N", this.Name, "T", t.Name),
		).
		Code(
			Switch(ˇ.This).CasesFrom(true, GEN_FOR(namesAndRenames, 2, func(namerename []Any) ISyn { // switch this
				return Case(N(namerename[0].(string)), // case ‹enumerantIdent›:
					ˇ.R.Set(L(namerename[1].(string))), // r = "‹enumerantNameOrRenamed›"
				)
			})...).DefaultCase( // default:
				GEN_BYCASE(USUALLY(ˇ.R.Set( // r =
					pkgstrconv.C("FormatInt", T.Int64.Conv(ˇ.This), L(10))), // pkg__strconv.FormatInt(int64(this), 10)
				), UNLESS{
					ebt == "int": ˇ.R.Set( // r =
						pkgstrconv.C("Itoa", T.Int.Conv(ˇ.This))), // pkg__strconv.Itoa(int(this))
					ustr.In(ebt, "uint", "uint8", "uint16", "uint32", "uint64"): ˇ.R.Set( // r =
						pkgstrconv.C("FormatUint", T.Uint64.Conv(ˇ.This), L(10))), // pkg__strconv.FormatUint(uint64(this), 10)
				}),
			),
		)
}

func (this *StringMethodOpts) genParseFunc(t *gent.Type, docComment gent.Str, minLen int, maybeCommonPrefix string, switchCases SynCases, defCase ...ISyn) *SynFunc {
	casehint, funcname := "and case-sensitively", this.Parser.FuncName.With("T", t.Name, "str", this.Name)
	if this.Parser.WithIgnoreCaseCmp {
		casehint = "but case-insensitively"
	}
	earlycheck := B.Len.C(ˇ.S).Lt(minLen) // len(s) < ‹minLen›}
	if l := len(maybeCommonPrefix); l > 0 && !this.Parser.WithIgnoreCaseCmp {
		earlycheck = earlycheck.Or(ˇ.S.Sl(0, l).Neq(maybeCommonPrefix)) // || s[0:5] != "PREF_" (for example)
	}
	return Func(funcname).Args(ˇ.S.OfType(T.String)).Rets(t.G.This, ˇ.Err).
		Doc(
			docComment.With("N", funcname, "T", t.Name, "s", ˇ.S.Name, "str", this.Name, "caseSensitivity", casehint),
		).
		Code(
			If(earlycheck, Then(
				GoTo("tryParseNum"),
			)),
			Switch(nil).
				CasesOf(switchCases...).
				DefaultCase(GoTo("tryParseNum")),
			K.Return,
			Label("tryParseNum", defCase...),
		)
}

func (this *GentStringersMethods) genParseErrlessFunc(t *gent.Type, funcName string, parseFuncName string) *SynFunc {
	maybe, fallback := N("maybe"+t.Name), N("fallback")
	return Func(funcName).Args(ˇ.S.OfType(T.String), fallback.OfType(t.G.T)).Rets(t.G.This).
		Doc(
			this.DocComments.ParsersErrlessVariant.With("N", funcName, "T", t.Name, "p", parseFuncName, "fallback", fallback.Name),
		).
		Code(
			Tup(maybe, ˇ.Err).Let(N(parseFuncName).C(ˇ.S)), // maybe,err := ‹parseFunc›(s)
			If(ˇ.Err.Eq(B.Nil), // if err == nil
				Then(ˇ.This.Set(maybe)),     // this = maybe
				Else(ˇ.This.Set(fallback))), // else this = fallback
		)
}

func (this *StringMethodOpts) genStringer(t *gent.Type, pkgstrconv PkgName) *SynFunc {
	ren, namesrenames := this.EnumerantRename != nil, make([]Any, 0, 2*len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if ren {
				renamed = this.EnumerantRename(enumerant)
			}
			namesrenames = append(namesrenames, enumerant, renamed)
		}
	}

	return this.genStringerMethod(t, pkgstrconv, namesrenames)
}

func (this *StringMethodOpts) genParser(t *gent.Type, docComment gent.Str, pkgstrconv PkgName, pkgstrings PkgName) *SynFunc {
	enstrs, switchcases := make([]string, 0, len(t.Enumish.ConstNames)), make(SynCases, 0, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if enstr := enumerant; enumerant != "_" {
			if rename := this.EnumerantRename; rename != nil {
				enstr = rename(enumerant)
			}
			enstrs = append(enstrs, enstr)
			enlit := L(enstr)
			var cmp IExprBoolish = ˇ.S.Eq(enlit)
			if this.Parser.WithIgnoreCaseCmp {
				cmp = cmp.Or(pkgstrings.C("EqualFold", ˇ.S, enlit))
			}
			switchcases.Add(cmp,
				ˇ.This.Set(N(enumerant)),
			)
		}
	}

	enumbasetype, defaultcase := TrNamed("", t.Enumish.BaseType), func(ebt *TypeRef, parse string, args ...ISyn) Syns {
		vtmp := N(ˇ.This.Name + t.Enumish.BaseType)
		return Syns{Var(vtmp.Name, ebt, nil),
			Tup(vtmp, ˇ.Err).Set(pkgstrconv.C(parse, append(Syns{ˇ.S}, args...)...)),
			If(ˇ.Err.Eq(B.Nil), Then(
				ˇ.This.Set(t.G.T.Conv(vtmp)),
			)),
		}
	}
	var switchdefault Syns
	switch enumbasetype.Named.TypeName {
	case "int":
		switchdefault = defaultcase(T.Int, "Atoi")
	case "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		switchdefault = defaultcase(T.Uint64, "ParseUint", L(10), L(enumbasetype.BitSizeIfBuiltInNumberType()))
	case "int8", "int16", "int32", "int64", "rune":
		switchdefault = defaultcase(T.Int64, "ParseInt", L(10), L(enumbasetype.BitSizeIfBuiltInNumberType()))
	}

	return this.genParseFunc(t, docComment, len(ustr.Shortest(enstrs)), ustr.CommonPrefix(enstrs...), switchcases, switchdefault...)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStringersMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if len(this.Stringers) > 0 && t.IsEnumish() {
		yield = make(Syns, 0, 3*len(this.Stringers))
		pkgstrconv, pkgstrings := ctx.I("strconv"), ctx.I("strings")
		for i := range this.Stringers {
			if self := &this.Stringers[i]; !self.Disabled {
				if yield.Add(self.genStringer(t, pkgstrconv)); self.Parser.Add {
					fnp := self.genParser(t, this.DocComments.Parsers, pkgstrconv, pkgstrings)
					if yield.Add(fnp); self.Parser.Errless.Add {
						yield.Add(this.genParseErrlessFunc(t, fnp.Name+self.Parser.Errless.NameOrSuffix, fnp.Name))
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
	for i := range this.Stringers {
		this.Stringers[i].Disabled = disabled
		this.Stringers[i].Parser.Add = enabled
		this.Stringers[i].Parser.Errless.Add = enabled
	}
}
