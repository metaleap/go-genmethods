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

// With merges all `IGent`s in this with all those in `gents` into `merged`.
func (this Gents) With(gents ...Gents) (merged Gents) {
	merged = append(make(Gents, 0, len(this)+2*len(gents)), this...)
	for i := range gents {
		merged = append(merged, gents[i]...)
	}
	return
}

// EnableOrDisableAll sets all `IGent.Opt().Disabled` fields to `!enabled`.
func (this Gents) EnableOrDisableAll(enabled bool) {
	disabled := !enabled
	for i := range this {
		this[i].Opt().Disabled = disabled
	}
}

// EnableOrDisableAllVariantsAndOptionals calls
// the same-named method on all `IGent`s in `this`.
func (this Gents) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	for i := range this {
		this[i].EnableOrDisableAllVariantsAndOptionals(enabled)
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
func (this *Opts) EnableOrDisableAllVariantsAndOptionals(bool) {}

func (this *Opts) mayRunForType(ctx *Ctx, t *Type) bool {
	if this.Disabled || this.RunOnlyOnceWithoutAnyType || (this.RunOnlyForTypeAliases != t.Alias) {
		return false
	}
	if len(this.RunNeverForTypes.Named) > 0 {
		for _, tname := range this.RunNeverForTypes.Named {
			if tname == t.Name {
				return false
			}
		}
	}
	if this.RunNeverForTypes.Satisfying != nil && this.RunNeverForTypes.Satisfying(ctx, t) {
		return false
	}
	if len(this.RunOnlyForTypes.Named) > 0 {
		for _, tname := range this.RunOnlyForTypes.Named {
			if tname == t.Name {
				return true
			}
		}
		return false
	}
	if this.RunOnlyForTypes.Satisfying != nil {
		return this.RunOnlyForTypes.Satisfying(ctx, t)
	}

	return this.MayRunForType == nil || this.MayRunForType(ctx, t)
}

// Opt implements `IGent.Opt()` for `Opts` embedders.
func (this *Opts) Opt() *Opts { return this }

// RunGents instructs the given `gents` to generate code for `this` `Pkg`.
func (this *Pkg) RunGents(maybeCtxOpts *CtxOpts, gents Gents) (src []byte, stats *Stats, err error) {
	ctx, dst, codegencommentnotice :=
		maybeCtxOpts.newCtx(this, gents), udevgogen.File(this.Name, 2*len(this.Types)*len(gents)), CodeGenCommentNotice.With("progName", CodeGenCommentProgName)
	for _, g := range gents {
		if g.Opt().RunOnlyOnceWithoutAnyType {
			dst.Body.Add(ctx.generateTopLevelDecls(g, nil)...)
		}
	}
	for _, t := range this.Types {
		for _, g := range gents {
			if ctx.mayGentRunForType(g, t) {
				dst.Body.Add(ctx.generateTopLevelDecls(g, t)...)
			}
		}
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

// MustRunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `this`.
func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, statsPerPkg map[*Pkg]*Stats) {
	var errs map[*Pkg]error
	timeTakenTotal, statsPerPkg, errs = this.RunGentsAndGenerateOutputFiles(maybeCtxOpts, gents)
	for _, err := range errs {
		panic(err)
	}
	return
}

// RunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `this`.
func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, statsPerPkg map[*Pkg]*Stats, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup
	starttime, run := time.Now(), func(pkg *Pkg) {
		src, stats, err := pkg.RunGents(maybeCtxOpts, gents)
		if err == nil {
			err = ufs.WriteBinaryFile(filepath.Join(pkg.DirPath, pkg.CodeGen.OutputFileName), src)
		} else {
			println(string(src))
		}
		maps.Lock()
		if statsPerPkg[pkg] = stats; err != nil {
			errs[pkg] = err
		}
		maps.Unlock()
		runs.Done()
	}

	statsPerPkg, errs = make(map[*Pkg]*Stats, len(this)), map[*Pkg]error{}
	runs.Add(len(this))
	for _, pkg := range this {
		go run(pkg)
	}
	runs.Wait()
	timeTakenTotal = time.Since(starttime)
	return
}
