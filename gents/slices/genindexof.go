package gentslices

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

type GentIndexOfMethods struct {
	gent.Opts

	IndexMethod     IndexOfMethod
	LastIndexMethod IndexOfMethod
	IndexAnyMethod  IndexOfMethod
	IndicesMethod   IndexOfMethod
}

type IndexOfMethod struct {
	Disabled          bool
	DocComment        gent.Str
	Name              string
	FuncVariationName string
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (this *GentIndexOfMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.Ast.TArrOrSl != nil {
		println("NOICE")
	}
	return
}
