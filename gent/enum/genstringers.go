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
	caseof, pkgstrconv := Switch(V.This), N(t.Pkg.I("strconv"))
	caseof.Cases = make([]SynCond, 0, len(t.Enumish.ConstNames))
	for _, enumerant := range t.Enumish.ConstNames {
		if renamed := enumerant; enumerant != "_" {
			if rename := this.Stringers[strName]; rename != nil {
				renamed = rename(renamed)
			}
			caseof.Cases = append(caseof.Cases, Cond(N(enumerant), Set(V.Ret, L(renamed))))
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
	fn = Fn(NoMethodRecv, strings.NewReplacer("{T}", t.Name, "{s}", strName).Replace(this.Parsers.FuncName), TdFunc(NTs("s", T.String), t.CodeGen.MethodRecvVal, V.Err),
		K.Ret,
	)
	return
}
