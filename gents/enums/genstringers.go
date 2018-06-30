package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
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

func (this *GentStringersMethods) genStringerMethod(self *StringMethodOpts, t *gent.Type, pkgstrconv PkgName) *SynFunc {
	switchcase := Switch(V.This, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := self.EnumerantRename; rename != nil {
				renamed = rename(renamed)
			}
			switchcase.Case(N(enumerant),
				V.R.SetTo(L(renamed)))
		}
	}

	switch t.Enumish.BaseType {
	case "int":
		switchcase.DefaultCase(
			V.R.SetTo(pkgstrconv.C("Itoa", T.Int.Conv(V.This))))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		switchcase.DefaultCase(
			V.R.SetTo(pkgstrconv.C("FormatUint", T.Uint64.Conv(V.This), L(10))))
	default:
		switchcase.DefaultCase(
			V.R.SetTo(pkgstrconv.C("FormatInt", T.Int64.Conv(V.This), L(10))))
	}
	return t.G.ThisVal.Method(self.Name).Sig(&Sigs.NoneToString).
		Doc(
			self.DocComment.With("N", self.Name, "T", t.Name),
		).
		Code(
			switchcase,
		)
}

func (this *GentStringersMethods) genParseFunc(self *StringMethodOpts, t *gent.Type, pkgstrconv PkgName, pkgstrings PkgName) *SynFunc {
	s, switchcase := N("s"), Switch(nil, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := L(enumerant); enumerant != "_" {
			if rename := self.EnumerantRename; rename != nil {
				renamed = L(rename(enumerant))
			}

			var cmp IDotsBoolish = s.Eq(renamed)
			if self.ParseAddIgnoreCaseCmp {
				cmp = cmp.Or(pkgstrings.C("EqualFold", s, renamed))
			}
			switchcase.Case(cmp,
				V.This.SetTo(N(enumerant)),
			)
		}
	}

	vtmp, enumbasetype := N(V.This.Name+t.Enumish.BaseType), TrNamed("", t.Enumish.BaseType)
	adddefault := func(tref *TypeRef, strconvparsefunc string, args ...ISyn) {
		switchcase.DefaultCase(
			Var(vtmp.Name, tref, nil),
			Tup(vtmp, V.Err).SetTo(pkgstrconv.C(strconvparsefunc, append(Syns{s}, args...)...)),
			IfThen(V.Err.Eq(B.Nil),
				V.This.SetTo(t.G.T.Conv(vtmp)),
			),
		)
	}
	switch enumbasetype.Named.TypeName {
	case "int":
		adddefault(T.Int, "Atoi")
	case "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		adddefault(T.Uint64, "ParseUint", L(10), L(enumbasetype.SafeBitSizeIfBuiltInNumberType()))
	case "int8", "int16", "int32", "int64", "rune":
		adddefault(T.Int64, "ParseInt", L(10), L(enumbasetype.SafeBitSizeIfBuiltInNumberType()))
	}

	casehint, funcname := "and case-sensitively", self.ParseFuncName.With("T", t.Name, "str", self.Name)
	if self.ParseAddIgnoreCaseCmp {
		casehint = "but case-insensitively"
	}
	return Func(funcname).Args(s.T(T.String)).Rets(t.G.ThisVal, V.Err).
		Doc(
			this.DocComments.Parsers.With("N", funcname, "T", t.Name, "s", s.Name, "str", self.Name, "caseSensitivity", casehint),
		).
		Code(
			switchcase,
		)
}

func (this *GentStringersMethods) genParseErrlessFunc(t *gent.Type, funcName string, parseFuncName string) *SynFunc {
	s, maybe, fallback := N("s"), N("maybe"+t.Name), N("fallback")
	return Func(funcName).Arg(s.Name, T.String).Arg(fallback.Name, t.G.ThisVal.Type).Rets(t.G.ThisVal).
		Doc(
			this.DocComments.ParsersErrlessVariant.With("N", funcName, "T", t.Name, "p", parseFuncName, "fallback", fallback.Name),
		).
		Code(
			Tup(maybe, V.Err).Decl(C.Named(parseFuncName, s)),
			If(V.Err.Eq(B.Nil),
				Then(V.This.SetTo(maybe)),
				Else(V.This.SetTo(fallback))),
		)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStringersMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if len(this.Stringers) > 0 && t.IsEnumish() {
		decls = make(Syns, 0, 2+len(t.Enumish.ConstNames)*3*len(this.Stringers))
		pkgstrconv, pkgstrings := ctx.I("strconv"), ctx.I("strings")
		for i := range this.Stringers {
			if self := &this.Stringers[i]; !self.Disabled {
				decls.Add(this.genStringerMethod(self, t, pkgstrconv))
				if self.ParseFuncName != "" {
					fnp := this.genParseFunc(self, t, pkgstrconv, pkgstrings)
					if decls.Add(fnp); self.ParseErrless.Add {
						decls.Add(this.genParseErrlessFunc(t, fnp.Name+self.ParseErrless.NameOrSuffix, fnp.Name))
					}
				}
			}
		}
	}
	return
}
