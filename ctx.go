package gent

import (
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/sys"
)

func init() {
	Defaults.Ctx.Opt.EmitNoOpFuncBodies = usys.EnvBool("GOGENT_EMITNOOPS", false)
	Defaults.Ctx.Opt.NoGoFmt = usys.EnvBool("GOGENT_NOGOFMT", false)
}

type ctxDeclKey struct {
	g IGent
	t *Type
}

type Ctx struct {
	// Opt holds the only user-settable fields (in between runs).
	// Code-gens only read but don't mutate them.
	Opt struct {
		NoGoFmt            bool
		EmitNoOpFuncBodies bool
	}

	TimeStarted time.Time

	pkgImportPathsToPkgImportNames udevgogen.PkgImports

	declsGenerated map[ctxDeclKey]udevgogen.Syns
}

func (this *Ctx) anew() *Ctx {
	if this == nil {
		this = &Defaults.Ctx
	}
	return &Ctx{
		Opt: this.Opt, TimeStarted: time.Now(), pkgImportPathsToPkgImportNames: udevgogen.PkgImports{},
		declsGenerated: map[ctxDeclKey]udevgogen.Syns{},
	}
}

func (this *Ctx) generateTopLevelDecls(g IGent, t *Type) (decls udevgogen.Syns) {
	decls = g.GenerateTopLevelDecls(this, t)
	this.declsGenerated[ctxDeclKey{g: g, t: t}] = decls
	return
}

func (this *Ctx) I(pkgImportPath string) (pkgImportName string) {
	pkgImportName = this.pkgImportPathsToPkgImportNames.Ensure(pkgImportPath)
	return
}

func (this *Ctx) DeclsGeneratedSoFar(maybeGent IGent, maybeType *Type) (matches []udevgogen.Syns) {
	for key, decls := range this.declsGenerated {
		if (maybeGent == nil || key.g == maybeGent) && (maybeType == nil || key.t == maybeType) {
			matches = append(matches, decls)
		}
	}
	return
}
