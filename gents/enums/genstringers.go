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

func (this *GentStringersMethods) genStringerMethod(idx int, ctx *gent.Ctx, t *gent.Type, pkgstrconv PkgName) *SynFunc {
	self, switchcase := &this.Stringers[idx], Switch(V.This, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := self.EnumerantRename; rename != nil {
				renamed = rename(renamed)
			}
			switchcase.Cases.Add(N(enumerant),
				V.R.SetTo(L(renamed)))
		}
	}

	var defcase OpSet
	switch t.Enumish.BaseType {
	case "int":
		defcase = V.R.SetTo(pkgstrconv.C("Itoa", T.Int.Conv(V.This)))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		defcase = V.R.SetTo(pkgstrconv.C("FormatUint", T.Uint64.Conv(V.This), L(10)))
	default:
		defcase = V.R.SetTo(pkgstrconv.C("FormatInt", T.Int64.Conv(V.This), L(10)))
	}
	switchcase.Default.Add(defcase)
	return t.Gen.ThisVal.Method(self.Name).Sig(&Sigs.NoneToString).
		Doc(self.DocComment.With("N", self.Name, "T", t.Name)).
		Code(switchcase)
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStringersMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if len(this.Stringers) > 0 && t.IsEnumish() {
		decls = make(Syns, 0, 2+len(t.Enumish.ConstNames)*3*len(this.Stringers))
		pkgstrconv := ctx.I("strconv")
		for i := range this.Stringers {
			if !this.Stringers[i].Disabled {
				decls.Add(this.genStringerMethod(i, ctx, t, pkgstrconv))
				if this.Stringers[i].ParseFuncName != "" {
					decls.Add(this.genParser(i, ctx, t)...)
				}
			}
		}
	}
	return
}

func (this *GentStringersMethods) genParser(idx int, ctx *gent.Ctx, t *gent.Type) (synFuncs Syns) {
	self, s, caseof, pkgstrconv := &this.Stringers[idx], N("s"), Switch(nil, len(t.Enumish.ConstNames)), ctx.I("strconv")
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := L(enumerant); enumerant != "_" {
			if rename := self.EnumerantRename; rename != nil {
				renamed = L(rename(enumerant))
			}
			var cmp ISyn = Eq(s, renamed)
			if self.ParseAddIgnoreCaseCmp {
				cmp = Or(cmp, ctx.I("strings").C("EqualFold", s, renamed))
			}
			caseof.Cases.Add(cmp, Set(V.This, N(enumerant)))
		}
	}

	vn, enumbasetype := N(V.This.Name+t.Enumish.BaseType), TrNamed("", t.Enumish.BaseType)
	adddefault := func(tref *TypeRef, callName string, callArgs ...ISyn) {
		caseof.Default.Add(
			Var(vn.Name, tref, nil),
			Set(Tup(vn, V.Err), pkgstrconv.C(callName, append(Syns{s}, callArgs...)...)),
			If(Eq(V.Err, B.Nil), Set(V.This, C.N(t.Name, vn))),
		)
	}
	switch t.Enumish.BaseType {
	case "int":
		adddefault(T.Int, "Atoi")
	case "uint", "uint8", "uint16", "uint32", "uint64", "byte":
		adddefault(T.Uint64, "ParseUint", L(10), L(enumbasetype.SafeBitSizeIfBuiltInNumberType()))
	case "int8", "int16", "int32", "int64", "rune":
		adddefault(T.Int64, "ParseInt", L(10), L(enumbasetype.SafeBitSizeIfBuiltInNumberType()))
	}

	fname := self.ParseFuncName.With("T", t.Name, "str", self.Name)
	fnp := Fn(NoMethodRecv, fname, TdFunc(NTs(s.Name, T.String), t.Gen.ThisVal, V.Err),
		caseof,
	)
	doccs := "and case-sensitively"
	if self.ParseAddIgnoreCaseCmp {
		doccs = "but case-insensitively"
	}
	fnp.Docs.Add(this.DocComments.Parsers.With("N", fnp.Name, "T", t.Name, "s", s.Name, "str", self.Name, "caseSensitivity", doccs))
	synFuncs = Syns{fnp}

	if self.ParseErrless.Add {
		maybe, fallback := N("maybe"+t.Name), N("fallback")
		fnv := Fn(NoMethodRecv, fname+self.ParseErrless.NameOrSuffix, TdFunc(NTs(s.Name, T.String, fallback.Name, t.Gen.ThisVal.Type), t.Gen.ThisVal),
			Decl(Tup(maybe, V.Err.Named), C.N(fname, s)),
			Ifs(Eq(V.Err, B.Nil),
				Block(Set(V.This, maybe)),
				Block(Set(V.This, N("fallback")))),
		)
		fnv.Docs.Add(this.DocComments.ParsersErrlessVariant.With("N", fnv.Name, "T", t.Name, "p", fnp.Name, "fallback", fallback.Name))
		synFuncs = append(synFuncs, fnv)
	}
	return
}
