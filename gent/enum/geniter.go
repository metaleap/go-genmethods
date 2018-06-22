package gentenum

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentIterateFunc struct {
	FuncName       string
	EnumerantName  bool
	EnumerantValue bool
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentIterateFunc) GenerateTopLevelDecls(t *gent.Type) (tlDecls Syns) {
	if t.SeemsEnumish() {

	}
	return
}
