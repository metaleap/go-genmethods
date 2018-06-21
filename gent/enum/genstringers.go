package gentenum

import (
	"strings"

	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentStringMethods struct {
	DocComment string
	Stringers  map[string]func(string) string
	Parsers    struct {
		OnePerStringer         bool
		OneUber                bool
		AddIgnoreCaseCmp       bool
		FuncName               string
		AddErrlessWithFallback bool
	}
}

func (this *GentStringMethods) GenerateTopLevelDecls(t *gent.Type) (tlDecls Syns) {
	if len(this.Stringers) > 0 && t.SeemsEnumish() {
		tlDecls = make(Syns, 0, 2+len(t.Enumish.ConstNames)*3*len(this.Stringers))
		for strname := range this.Stringers {
			tlDecls.Add(this.genStringer(strname, t))
			if this.Parsers.OnePerStringer {
				tlDecls.Add(this.genParser(strname, t))
			}
		}
		if this.Parsers.OneUber {
		}
	}
	return
}

func (this *GentStringMethods) genStringer(strName string, t *gent.Type) (method *SynFunc) {
	caseof, pkgstrconv := Switch(V.This, len(t.Enumish.ConstNames)), N(t.Pkg.I("strconv"))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := this.Stringers[strName]; rename != nil {
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

	method = Fn(t.CodeGen.MethodRecvVal, strName, &Sigs.NoneToString,
		caseof,
	)
	return
}

func (this *GentStringMethods) genParser(strName string, t *gent.Type) (fn *SynFunc) {
	s, caseof, pkgstrconv := N("s"), Switch(nil, len(t.Enumish.ConstNames)), N(t.Pkg.I("strconv"))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := L(enumerant); enumerant != "_" {
			if rename := this.Stringers[strName]; rename != nil {
				renamed = L(rename(enumerant))
			}
			var cmp ISyn = Eq(s, renamed)
			if this.Parsers.AddIgnoreCaseCmp {
				cmp = Or(cmp, Call(D(N(t.Pkg.I("strings")), N("EqualFold")), s, renamed))
			}
			caseof.Cases.Add(cmp, Set(V.This, N(enumerant)))
		}
	}

	vn := N(V.This.Name + t.Enumish.BaseType)
	switch t.Enumish.BaseType {
	case "int":
		caseof.Default.Add(
			Var(vn.Name, T.Int, nil),
			Set(C(vn, V.Err), Call(D(pkgstrconv, N("Atoi")), s)),
			If(Eq(V.Err, B.Nil), Set(V.This, Call(N(t.Name), vn))),
		)
		// case "uint", "uint8", "uint16", "uint32", "uint64":
		// 	caseof.Default.Add(Set(V.Ret, Call(D(pkgstrconv, N("FormatUint")), Call(N("uint64"), V.This), L(10))))
		// default:
		// 	caseof.Default.Add(Set(V.Ret, Call(D(pkgstrconv, N("FormatInt")), Call(N("int64"), V.This), L(10))))
	}

	fn = Fn(NoMethodRecv, strings.NewReplacer("{T}", t.Name, "{s}", strName).Replace(this.Parsers.FuncName), TdFunc(NTs(s.Name, T.String), t.CodeGen.MethodRecvVal, V.Err),
		caseof,
	)
	return
}
