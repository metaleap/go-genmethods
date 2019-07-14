package gent

import (
	"strconv"
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/sys"
)

func init() {
	Defaults.CtxOpt.EmitNoOpFuncBodies = usys.EnvBool("GOGENT_EMITNOOPS", false)
	Defaults.CtxOpt.NoGoFmt = usys.EnvBool("GOGENT_NOGOFMT", false)
}

// CtxOpts wraps `Ctx` options.
type CtxOpts struct {
	// For Defaults.CtxOpts, initialized from env-var
	// `GOGENT_NOGOFMT` if `strconv.ParseBool`able.
	NoGoFmt bool

	// For Defaults.CtxOpts, initialized from env-var
	// `GOGENT_EMITNOOPS` if `strconv.ParseBool`able.
	EmitNoOpFuncBodies bool

	// If set, can be used to prevent running of the given
	// (or any) `IGent` on the given (or any) `*Type`.
	// See also `IGent.Opt().MayRunForType`.
	MayGentRunForType func(IGent, *Type) bool

	allTypeNames map[string]bool
}

type Stats struct {
	DurationOf struct {
		Constructing time.Duration
		Emitting     time.Duration
		Formatting   time.Duration
		Everything   time.Duration
	}
}

type ctxDeclKey struct {
	g IGent
	t *Type
}

// Ctx is a codegen-time context during a `Pkg.RunGents`
// call and is passed to `IGent.GenerateTopLevelDecls`.
type Ctx struct {
	// options pertaining to this `Ctx`
	Opt CtxOpts

	// strictly read-only
	Pkg *Pkg

	ExtraDefs []*udevgogen.SynFunc

	Gents                          Gents
	timeStarted                    time.Time
	declsGenerated                 map[ctxDeclKey]udevgogen.Syns
	pkgImportPathsToPkgImportNames udevgogen.PkgImports
	counter                        int
}

func (me *CtxOpts) newCtx(pkg *Pkg, gents Gents) *Ctx {
	if me == nil {
		me = &Defaults.CtxOpt
	}
	return &Ctx{
		Opt: *me, timeStarted: time.Now(), Gents: gents, Pkg: pkg,
		pkgImportPathsToPkgImportNames: udevgogen.PkgImports{},
		declsGenerated:                 map[ctxDeclKey]udevgogen.Syns{},
	}
}

func (me *Ctx) MayGentRunForType(g IGent, t *Type) bool {
	me.Opt.allTypeNames[t.Name] = true
	return g.Opt().mayRunForType(me, t) &&
		(me.Opt.MayGentRunForType == nil || me.Opt.MayGentRunForType(g, t))
}

func (me *Ctx) N(pref string) udevgogen.Named {
	me.counter++
	return udevgogen.N(pref + strconv.Itoa(me.counter))
}

func (me *Ctx) generateTopLevelDecls(g IGent, t *Type) (decls udevgogen.Syns) {
	decls = g.GenerateTopLevelDecls(me, t)
	me.declsGenerated[ctxDeclKey{g: g, t: t}] = decls
	if g.Opt().EmitCommented {
		for i := range decls {
			if fn, _ := decls[i].(*udevgogen.SynFunc); fn != nil {
				fn.EmitCommented = true
			} else if raw, _ := decls[i].(*udevgogen.SynRaw); raw != nil {
				raw.EmitCommented = true
			}
		}
	}
	return
}

// DeclsGeneratedSoFar collects and returns all results of `IGent.GenerateTopLevelDecls`
// performed so far by this `Ctx`, filtered optionally by `IGent` and/or by `Type`.
func (me *Ctx) DeclsGeneratedSoFar(maybeGent IGent, maybeType *Type) (matches []udevgogen.Syns) {
	for gt, decls := range me.declsGenerated {
		if (maybeGent == nil || gt.g == maybeGent) && (maybeType == nil || gt.t == maybeType) {
			matches = append(matches, decls)
		}
	}
	return
}

func (me *Ctx) GentExistsFor(t *Type, check func(IGent) bool) bool {
	if t != nil {
		for _, gent := range me.Gents {
			if me.MayGentRunForType(gent, t) && check(gent) {
				return true
			}
		}
	}
	return false
}

// Import returns the `pkgImportName` for the specified `pkgImportPath`.
// Eg. `Import("encoding/json")` might return `pkg__encoding_json` and
// more importantly, the import will be properly emitted (only if any of
// the import's uses get emitted) at code-gen time. Import is a `Ctx`-local
// wrapper of the `github.com/go-leap/dev/go/gen.PkgImports.Ensure` method.
func (me *Ctx) Import(pkgImportPath string) (pkgImportName udevgogen.PkgName) {
	pkgImportName = me.pkgImportPathsToPkgImportNames.Ensure(pkgImportPath)
	return
}
