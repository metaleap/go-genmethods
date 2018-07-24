package genttrav

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

const (
	DefaultMethodName = "TraverseFields"
)

func init() {
	Gents.StructFieldsTrav.MethodName = DefaultMethodName
}

type GentStructFieldsTrav struct {
	gent.Opts

	DocComment gent.Str
	MethodName string
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentStructFieldsTrav) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if tstruct := t.Expr.GenRef.Struct; tstruct != nil {
		// yield.Add(t.G.Tª.Method(this.MethodName))
	}
	return
}
