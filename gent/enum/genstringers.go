package gentenum

import (
	"strings"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentStringMethods struct {
	DocComment string
	Stringers  []Stringer
}

type Stringer struct {
	Name                   string
	EnumerantRename        func(string) string
	ParseFuncName          string
	ParseAddIgnoreCaseCmp  bool
	ParseAddErrlessVariant string
}

func (this *GentStringMethods) GenerateTopLevelDecls(t *gent.Type) (tlDecls Syns) {
	if len(this.Stringers) > 0 && t.SeemsEnumish() {
		tlDecls = make(Syns, 0, 2+len(t.Enumish.ConstNames)*3*len(this.Stringers))
		for i := range this.Stringers {
			tlDecls.Add(this.genStringer(i, t))
			if this.Stringers[i].ParseFuncName != "" {
				tlDecls.Add(this.genParser(i, t)...)
			}
		}
	}
	return
}

func (this *GentStringMethods) genStringer(idx int, t *gent.Type) (method *SynFunc) {
	caseof, pkgstrconv := Switch(V.This, len(t.Enumish.ConstNames)), N(t.Pkg.I("strconv"))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := this.Stringers[idx].EnumerantRename; rename != nil {
				renamed = rename(renamed)
			}
			caseof.Cases.Add(N(enumerant), Set(V.Ret, L(renamed)))
		}
	}

	switch t.Enumish.BaseType {
	case "int":
		caseof.Default.Add(Set(V.Ret, Call(D(pkgstrconv, N("Itoa")), Call(N("int"), V.This))))
	case "uint", "uint8", "uint16", "uint32", "uint64":
		caseof.Default.Add(Set(V.Ret, Call(D(pkgstrconv, N("FormatUint")), Call(N("uint64"), V.This), L(10))))
	default:
		caseof.Default.Add(Set(V.Ret, Call(D(pkgstrconv, N("FormatInt")), Call(N("int64"), V.This), L(10))))
	}

	method = Fn(t.CodeGen.ThisVal, this.Stringers[idx].Name, &Sigs.NoneToString,
		caseof,
	)
	return
}

func (this *GentStringMethods) genParser(idx int, t *gent.Type) (synFuncs Syns) {
	s, caseof, pkgstrconv := N("s"), Switch(nil, len(t.Enumish.ConstNames)), N(t.Pkg.I("strconv"))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := L(enumerant); enumerant != "_" {
			if rename := this.Stringers[idx].EnumerantRename; rename != nil {
				renamed = L(rename(enumerant))
			}
			var cmp ISyn = Eq(s, renamed)
			if this.Stringers[idx].ParseAddIgnoreCaseCmp {
				cmp = Or(cmp, Call(D(N(t.Pkg.I("strings")), N("EqualFold")), s, renamed))
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

	fname := strings.NewReplacer("{T}", t.Name, "{s}", this.Stringers[idx].Name).Replace(this.Stringers[idx].ParseFuncName)
	synFuncs = Syns{Fn(NoMethodRecv, fname, TdFunc(NTs(s.Name, T.String), t.CodeGen.ThisVal, V.Err),
		caseof,
	)}
	if fnamesuffix := this.Stringers[idx].ParseAddErrlessVariant; fnamesuffix != "" {
		maybe := N("maybe" + t.Name)
		synFuncs = append(synFuncs, Fn(NoMethodRecv, fname+fnamesuffix, TdFunc(NTs(s.Name, T.String, "fallback", t.CodeGen.ThisVal.Type), t.CodeGen.ThisVal),
			Decl(C(maybe, V.Err.Named), Call(N(fname), s)),
			Ifs(Eq(V.Err, B.Nil),
				Block(Set(V.This, maybe)),
				Block(Set(V.This, N("fallback")))),
		))
	}
	return
}
