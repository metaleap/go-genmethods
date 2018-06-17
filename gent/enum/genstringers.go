package gentenum

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentStringMethods struct {
	DocComment string
	Stringers  map[string]func(string) string
	Parsers    struct {
		OnePerStringer         bool
		OneUber                bool
		AddErrlessWithFallback bool
	}
}

func (this *GentStringMethods) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn) {
	if len(this.Stringers) > 0 && t.Enumish.Potentially && len(t.Enumish.ConstNames) > 0 {

	}
	return
}
