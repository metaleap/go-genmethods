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
		{DocComment: DefaultStringers0DocComments, Name: DefaultStringers0MethodName,
			EnumerantRename: nil, ParseFuncName: DefaultStringersParsersFuncName, ParseErrless: gent.Variant{Add: false, NameOrSuffix: "Or"}},
		{DocComment: DefaultStringers1DocComments, Name: DefaultStringers1MethodName,
			Disabled: true, ParseFuncName: DefaultStringersParsersFuncName, ParseErrless: gent.Variant{Add: false, NameOrSuffix: "Or"}},
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
	Disabled              bool
	DocComment            gent.Str
	Name                  string
	EnumerantRename       func(string) string
	ParseFuncName         gent.Str
	ParseAddIgnoreCaseCmp bool
	ParseErrless          gent.Variant
}

func (this *StringMethodOpts) genStringerMethod(t *gent.Type, switchCases SynConds, defCase ...ISyn) *SynFunc {
	return t.G.ThisVal.Method(this.Name).Sig(&Sigs.NoneToString).
		Doc(
			this.DocComment.With("N", this.Name, "T", t.Name),
		).
		Code(
			Switch(ˇ.This, len(switchCases)+1).
				CasesOf(switchCases...).
				DefaultCase(defCase...),
		)
}

func (this *StringMethodOpts) genParseFunc(t *gent.Type, docComment gent.Str, minLen int, maybeCommonPrefix string, switchCases SynConds, defCase ...ISyn) *SynFunc {
	casehint, funcname := "and case-sensitively", this.ParseFuncName.With("T", t.Name, "str", this.Name)
	if this.ParseAddIgnoreCaseCmp {
		casehint = "but case-insensitively"
	}
	var earlycheck IExprBoolish = C.Len(ˇ.S).Lt(minLen) // len(s) < {minLen}
	if maybeCommonPrefix != "" {
		earlycheck = earlycheck.Or(ˇ.S.Sl(0, len(maybeCommonPrefix)).Neq(maybeCommonPrefix)) // || s[0:5] != "PREF_" (for example)
	}
	return Func(funcname).Args(ˇ.S.OfType(T.String)).Rets(t.G.ThisVal, ˇ.Err).
		Doc(
			docComment.With("N", funcname, "T", t.Name, "s", ˇ.S.Name, "str", this.Name, "caseSensitivity", casehint),
		).
		Code(
			If(earlycheck, Then(
				GoTo("tryParseNum"),
			)),
			Switch(nil, len(switchCases)+1).
				CasesOf(switchCases...).
				DefaultCase(GoTo("tryParseNum")),
			K.Return,
			Label("tryParseNum", defCase...),
		)
}

func (this *GentStringersMethods) genParseErrlessFunc(t *gent.Type, funcName string, parseFuncName string) *SynFunc {
	maybe, fallback := N("maybe"+t.Name), N("fallback")
	return Func(funcName).Arg(ˇ.S.Name, T.String).Arg(fallback.Name, t.G.ThisVal.Type).Rets(t.G.ThisVal).
		Doc(
			this.DocComments.ParsersErrlessVariant.With("N", funcName, "T", t.Name, "p", parseFuncName, "fallback", fallback.Name),
		).
		Code(
			Tup(maybe, ˇ.Err).Let(C.Named(parseFuncName, ˇ.S)),
			If(ˇ.Err.Eq(B.Nil),
				Then(ˇ.This.Set(maybe)),
				Else(ˇ.This.Set(fallback))),
		)
}

func (this *StringMethodOpts) genStringer(t *gent.Type, pkgstrconv PkgName) *SynFunc {
	switchcases := make(SynConds, 0, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := this.EnumerantRename; rename != nil {
				renamed = rename(renamed)
			}
			switchcases.Add(N(enumerant),
				ˇ.R.Set(L(renamed)))
		}
	}

	var switchdefault ISyn
	switch t.Enumish.BaseType {
	case "int":
		switchdefault = ˇ.R.Set(pkgstrconv.C("Itoa", T.Int.Conv(ˇ.This)))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		switchdefault = ˇ.R.Set(pkgstrconv.C("FormatUint", T.Uint64.Conv(ˇ.This), L(10)))
	default:
		switchdefault = ˇ.R.Set(pkgstrconv.C("FormatInt", T.Int64.Conv(ˇ.This), L(10)))
	}

	return this.genStringerMethod(t, switchcases, switchdefault)
}

func (this *StringMethodOpts) genParser(t *gent.Type, docComment gent.Str, pkgstrconv PkgName, pkgstrings PkgName) *SynFunc {
	enstrs, switchcases := make([]string, 0, len(t.Enumish.ConstNames)), make(SynConds, 0, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if enstr := enumerant; enumerant != "_" {
			if rename := this.EnumerantRename; rename != nil {
				enstr = rename(enumerant)
			}
			enstrs = append(enstrs, enstr)
			enlit := L(enstr)
			var cmp IExprBoolish = ˇ.S.Eq(enlit)
			if this.ParseAddIgnoreCaseCmp {
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
func (this *GentStringersMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if len(this.Stringers) > 0 && t.IsEnumish() {
		decls = make(Syns, 0, 2+len(t.Enumish.ConstNames)*3*len(this.Stringers))
		pkgstrconv, pkgstrings := ctx.I("strconv"), ctx.I("strings")
		for i := range this.Stringers {
			if self := &this.Stringers[i]; !self.Disabled {
				decls.Add(self.genStringer(t, pkgstrconv))
				if self.ParseFuncName != "" {
					fnp := self.genParser(t, this.DocComments.Parsers, pkgstrconv, pkgstrings)
					if decls.Add(fnp); self.ParseErrless.Add {
						decls.Add(this.genParseErrlessFunc(t, fnp.Name+self.ParseErrless.NameOrSuffix, fnp.Name))
					}
				}
			}
		}
	}
	return
}
