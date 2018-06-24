package gent

import (
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/sys"
)

func init() {
	Defaults.CtxOpt.EmitNoOpFuncBodies = usys.EnvBool("GOGENT_EMITNOOPS", false)
	Defaults.CtxOpt.NoGoFmt = usys.EnvBool("GOGENT_NOGOFMT", false)
}

type Opts struct {
	// For Defaults.CtxOpt, initialized from env-var
	// `GOGENT_NOGOFMT` if `strconv.ParseBool`able.
	NoGoFmt bool

	// For Defaults.CtxOpt, initialized from env-var
	// `GOGENT_EMITNOOPS` if `strconv.ParseBool`able.
	EmitNoOpFuncBodies bool

	// If set, can be used to prevent running of the given
	// (or any) `IGent` on the given (or any) `*Type`.
	MayGentRunForType func(IGent, *Type) bool
}

type ctxDeclKey struct {
	g IGent
	t *Type
}

type Ctx struct {
	Opt Opts

	timeStarted                    time.Time
	declsGenerated                 map[ctxDeclKey]udevgogen.Syns
	pkgImportPathsToPkgImportNames udevgogen.PkgImports
}

func (this *Opts) newCtx() *Ctx {
	if this == nil {
		this = &Defaults.CtxOpt
	}
	return &Ctx{
		Opt: *this, timeStarted: time.Now(), pkgImportPathsToPkgImportNames: udevgogen.PkgImports{},
		declsGenerated: map[ctxDeclKey]udevgogen.Syns{},
	}
}

func (this *Ctx) generateTopLevelDecls(g IGent, t *Type) (decls udevgogen.Syns) {
	decls = g.GenerateTopLevelDecls(this, t)
	this.declsGenerated[ctxDeclKey{g: g, t: t}] = decls
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

func (this *Ctx) I(pkgImportPath string) (pkgImportName string) {
	pkgImportName = this.pkgImportPathsToPkgImportNames.Ensure(pkgImportPath)
	return
}
