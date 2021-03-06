package gent

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/fs"
)

// IGent is the interface implemented by individual code-gens.
type IGent interface {
	// must never return `nil` (easiest impl is to embed `Opts`)
	Opt() *Opts

	// implemented as a no-op by `Opts`, to be
	// overridden by implementations as desired
	EnableOrDisableAllVariantsAndOptionals(bool)

	// may read from but never mutate its args.
	// expected to generate preferentially funcs / methods
	// instead of top-level const / var / type decls
	GenerateTopLevelDecls(*Ctx, *Type) udevgogen.Syns
}

// Gents is a slice if `IGent`s.
type Gents []IGent

// With merges all `IGent`s in `me` with all those in `gents` into `merged`.
func (me Gents) With(gents ...Gents) (merged Gents) {
	merged = append(make(Gents, 0, len(me)+2*len(gents)), me...)
	for i := range gents {
		merged = append(merged, gents[i]...)
	}
	return
}

// EnableOrDisableAll sets all `IGent.Opt().Disabled` fields to `!enabled`.
func (me Gents) EnableOrDisableAll(enabled bool) {
	disabled := !enabled
	for i := range me {
		me[i].Opt().Disabled = disabled
	}
}

// EnableOrDisableAllVariantsAndOptionals calls
// the same-named method on all `IGent`s in `me`.
func (me Gents) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	for i := range me {
		me[i].EnableOrDisableAllVariantsAndOptionals(enabled)
	}
}

// Opts related to a single `IGent`, and designed for embedding.
type Opts struct {
	Disabled      bool
	EmitCommented bool

	// not used by `go-gent`, but could be handy for callers
	Name string

	// tested right after `Disabled` and before the below checks.
	// should typically be false for most `Gent`s as the design assumes
	// generation of methods, but it's for the occasional need to generate
	// non-method `func`s related to certain type-alias declarations
	RunOnlyForTypeAliases bool

	RunOnlyOnceWithoutAnyType bool

	// if-and-only-if these are set, they're checked
	// before `MayRunForType` (but after `Disabled`)
	RunNeverForTypes, RunOnlyForTypes struct {
		Named      []string
		Satisfying func(*Ctx, *Type) bool
	}

	// If set, can be used to prevent running of
	// this `IGent` on the given (or any) `*Type`.
	// See also `CtxOpts.MayGentRunForType`.
	MayRunForType func(*Ctx, *Type) bool
}

// EnableOrDisableAllVariantsAndOptionals implements `IGent` but
// with a no-op, to be overridden by `Opts`-embedders as desired.
//
// To disable or enable an `IGent` itself, set `Opts.Disabled`.
func (me *Opts) EnableOrDisableAllVariantsAndOptionals(bool) {}

func (me *Opts) mayRunForType(ctx *Ctx, t *Type) bool {
	if me.Disabled || me.RunOnlyOnceWithoutAnyType || (me.RunOnlyForTypeAliases != t.Alias) {
		return false
	}
	if len(me.RunNeverForTypes.Named) > 0 {
		for _, tname := range me.RunNeverForTypes.Named {
			if tname == t.Name {
				return false
			}
		}
	}
	if me.RunNeverForTypes.Satisfying != nil && me.RunNeverForTypes.Satisfying(ctx, t) {
		return false
	}
	if len(me.RunOnlyForTypes.Named) > 0 {
		for _, tname := range me.RunOnlyForTypes.Named {
			if tname == t.Name {
				return true
			}
		}
		return false
	}
	if me.RunOnlyForTypes.Satisfying != nil {
		return me.RunOnlyForTypes.Satisfying(ctx, t)
	}

	return me.MayRunForType == nil || me.MayRunForType(ctx, t)
}

// Opt implements `IGent.Opt()` for `Opts` embedders.
func (me *Opts) Opt() *Opts { return me }

// RunGents instructs the given `gents` to generate code for `me`.
func (me *Pkg) RunGents(maybeCtxOpts *CtxOpts, gents Gents) (src []byte, stats *Stats, err error) {
	ctx, dst, codegencommentnotice :=
		maybeCtxOpts.newCtx(me, gents), udevgogen.File(me.Name, 2*len(me.Types)*len(gents)), CodeGenCommentNotice.With("progName", CodeGenCommentProgName)
	dst.DocComments = me.CodeGen.OutputFile.DocComments
	for _, g := range gents {
		if g.Opt().RunOnlyOnceWithoutAnyType {
			dst.Body.Add(ctx.generateTopLevelDecls(g, nil)...)
		}
	}
	for _, t := range me.Types {
		for _, g := range gents {
			if ctx.MayGentRunForType(g, t) {
				dst.Body.Add(ctx.generateTopLevelDecls(g, t)...)
			}
		}
	}
	for _, defcode := range ctx.ExtraDefs {
		dst.Add(defcode)
	}

	stats = &Stats{}
	stats.DurationOf.Constructing = time.Since(ctx.timeStarted)

	var timetakengofmt time.Duration
	emitstarttime := time.Now()
	src, timetakengofmt, err = dst.CodeGen(codegencommentnotice, ctx.pkgImportPathsToPkgImportNames, ctx.Opt.EmitNoOpFuncBodies, !ctx.Opt.NoGoFmt)
	stats.DurationOf.Formatting, stats.DurationOf.Emitting, stats.DurationOf.Everything =
		timetakengofmt, time.Since(emitstarttime)-timetakengofmt, time.Since(ctx.timeStarted)
	return
}

func (me *Pkg) RunGentsAndGenerateOutputFile(maybeCtxOpts *CtxOpts, gents Gents) (*Stats, error) {
	src, stats, err := me.RunGents(maybeCtxOpts, gents)
	if err == nil {
		err = ufs.WriteBinaryFile(filepath.Join(me.DirPath, me.CodeGen.OutputFile.Name), src)
	}
	return stats, err
}

// MustRunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `me`.
func (me Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, statsPerPkg map[*Pkg]*Stats) {
	var errs map[*Pkg]error
	timeTakenTotal, statsPerPkg, errs = me.RunGentsAndGenerateOutputFiles(maybeCtxOpts, gents)
	for _, err := range errs {
		panic(err)
	}
	return
}

// RunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `me`.
func (me Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, statsPerPkg map[*Pkg]*Stats, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup

	if maybeCtxOpts == nil {
		maybeCtxOpts = &Defaults.CtxOpt
	}
	if maybeCtxOpts.allTypeNames == nil {
		maybeCtxOpts.allTypeNames = make(map[string]bool, 32)
	}

	starttime := time.Now()
	run := func(pkg *Pkg) {
		stats, err := pkg.RunGentsAndGenerateOutputFile(maybeCtxOpts, gents)
		maps.Lock()
		if statsPerPkg[pkg] = stats; err != nil {
			errs[pkg] = err
		}
		maps.Unlock()
		runs.Done()
	}

	for _, gent := range gents {
		opt := gent.Opt()
		for _, tname := range opt.RunNeverForTypes.Named {
			maybeCtxOpts.allTypeNames[tname] = false
		}
		for _, tname := range opt.RunOnlyForTypes.Named {
			maybeCtxOpts.allTypeNames[tname] = false
		}
	}

	statsPerPkg, errs = make(map[*Pkg]*Stats, len(me)), map[*Pkg]error{}
	runs.Add(len(me))
	for _, pkg := range me {
		go run(pkg)
	}
	runs.Wait()
	timeTakenTotal = time.Since(starttime)

	for tname, tencountered := range maybeCtxOpts.allTypeNames {
		if !tencountered {
			println("Note: type name `" + tname + "` never encountered")
		}
	}

	return
}
