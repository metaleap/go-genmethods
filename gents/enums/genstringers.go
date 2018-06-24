package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentStringMethods generates for enum type-defs the specified
// `string`ifying methods, optionally with corresponding "parsing" funcs.
//
// An instance with illustrative defaults is in `Defaults.String`.
type GentStringMethods struct {
	gent.Opts

	Stringers   []StringMethod
	DocComments struct {
		Parsers               gent.Str
		ParsersErrlessVariant gent.Str
	}
}

type StringMethod struct {
	Disabled                         bool
	DocComment                       gent.Str
	Name                             string
	EnumerantRename                  func(string) string
	ParseFuncName                    gent.Str
	ParseAddIgnoreCaseCmp            bool
	ParseAddErrlessVariantWithSuffix string
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStringMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if len(this.Stringers) > 0 && t.IsEnumish() {
		decls = make(Syns, 0, 2+len(t.Enumish.ConstNames)*3*len(this.Stringers))
		for i := range this.Stringers {
			if !this.Stringers[i].Disabled {
				decls.Add(this.genStringer(i, ctx, t))
				if this.Stringers[i].ParseFuncName != "" {
					decls.Add(this.genParser(i, ctx, t)...)
				}
			}
		}
	}
	return
}

func (this *GentStringMethods) genStringer(idx int, ctx *gent.Ctx, t *gent.Type) (method *SynFunc) {
	self, caseof, pkgstrconv := &this.Stringers[idx], Switch(V.This, len(t.Enumish.ConstNames)), N(ctx.I("strconv"))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := self.EnumerantRename; rename != nil {
				renamed = rename(renamed)
			}
			caseof.Cases.Add(N(enumerant), Set(V.R, L(renamed)))
		}
	}

	switch t.Enumish.BaseType {
	case "int":
		caseof.Default.Add(Set(V.R, Call(D(pkgstrconv, N("Itoa")), Call(N("int"), V.This))))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		caseof.Default.Add(Set(V.R, Call(D(pkgstrconv, N("FormatUint")), Call(N("uint64"), V.This), L(10))))
	default:
		caseof.Default.Add(Set(V.R, Call(D(pkgstrconv, N("FormatInt")), Call(N("int64"), V.This), L(10))))
	}

	method = Fn(t.CodeGen.ThisVal, self.Name, &Sigs.NoneToString,
		caseof,
	)
	if self.DocComment != "" {
		method.Doc.Add(self.DocComment.With("{N}", method.Name, "{T}", t.Name))
	}
	return
}

func (this *GentStringMethods) genParser(idx int, ctx *gent.Ctx, t *gent.Type) (synFuncs Syns) {
	self, s, caseof, pkgstrconv := &this.Stringers[idx], N("s"), Switch(nil, len(t.Enumish.ConstNames)), N(ctx.I("strconv"))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := L(enumerant); enumerant != "_" {
			if rename := self.EnumerantRename; rename != nil {
				renamed = L(rename(enumerant))
			}
			var cmp ISyn = Eq(s, renamed)
			if self.ParseAddIgnoreCaseCmp {
				cmp = Or(cmp, Call(D(N(ctx.I("strings")), N("EqualFold")), s, renamed))
			}
			caseof.Cases.Add(cmp, Set(V.This, N(enumerant)))
		}
	}

	vn, enumbasetype := N(V.This.Name+t.Enumish.BaseType), TrNamed("", t.Enumish.BaseType)
	adddefault := func(tref *TypeRef, callName string, callArgs ...ISyn) {
		caseof.Default.Add(
			Var(vn.Name, tref, nil),
			Set(C(vn, V.Err), Call(D(pkgstrconv, N(callName)), append(Syns{s}, callArgs...)...)),
			If(Eq(V.Err, B.Nil), Set(V.This, Call(N(t.Name), vn))),
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

	fname := self.ParseFuncName.With("{T}", t.Name, "{str}", self.Name)
	fnp := Fn(NoMethodRecv, fname, TdFunc(NTs(s.Name, T.String), t.CodeGen.ThisVal, V.Err),
		caseof,
	)
	doccs := "and case-sensitively"
	if self.ParseAddIgnoreCaseCmp {
		doccs = "but case-insensitively"
	}
	fnp.Doc.Add(this.DocComments.Parsers.With("{N}", fnp.Name, "{T}", t.Name, "{s}", s.Name, "{str}", self.Name, "{caseSensitivity}", doccs))
	synFuncs = Syns{fnp}

	if fnvsuff := self.ParseAddErrlessVariantWithSuffix; fnvsuff != "" {
		maybe, fallback := N("maybe"+t.Name), N("fallback")
		fnv := Fn(NoMethodRecv, fname+fnvsuff, TdFunc(NTs(s.Name, T.String, fallback.Name, t.CodeGen.ThisVal.Type), t.CodeGen.ThisVal),
			Decl(C(maybe, V.Err.Named), Call(N(fname), s)),
			Ifs(Eq(V.Err, B.Nil),
				Block(Set(V.This, maybe)),
				Block(Set(V.This, N("fallback")))),
		)
		fnv.Doc.Add(this.DocComments.ParsersErrlessVariant.With("{N}", fnv.Name, "{T}", t.Name, "{p}", fnp.Name, "{fallback}", fallback.Name))
		synFuncs = append(synFuncs, fnv)
	}
	return
}
