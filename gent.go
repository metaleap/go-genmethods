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
	Disabled bool

	// not used by `go-gent`, but could be handy for callers
	Name string

	RunNeverForTypesNamed []string
	RunOnlyForTypesNamed  []string

	// If set, can be used to prevent running of
	// the `IGent` on the given (or any) `*Type`.
	// See also `CtxOpts.MayGentRunForType`.
	MayRunForType func(*Type) bool
}

// EnableOrDisableAllVariantsAndOptionals implements `IGent` but
// with a no-op, to be overridden by `Opts`-embedders as desired.
//
// To disable or enable an `IGent` itself, set `Opts.Disabled`.
func (this *Opts) EnableOrDisableAllVariantsAndOptionals(bool) {}

func (this *Opts) mayRunForType(t *Type) bool {
	if this.Disabled || t.Alias {
		return false
	}
	if len(this.RunNeverForTypesNamed) > 0 {
		for _, tname := range this.RunNeverForTypesNamed {
			if tname == t.Name {
				return false
			}
		}
	}
	if len(this.RunOnlyForTypesNamed) > 0 {
		for _, tname := range this.RunOnlyForTypesNamed {
			if tname == t.Name {
				return true
			}
		}
		return false
	}

	return this.MayRunForType == nil || this.MayRunForType(t)
}

// Opt implements `IGent.Opt()` for `Opts` embedders.
func (this *Opts) Opt() *Opts { return this }

// RunGents instructs the given `gents` to generate code for `this` `Pkg`.
func (this *Pkg) RunGents(maybeCtxOpts *CtxOpts, gents Gents) (src []byte, timetakengofmt time.Duration, err error) {
	ctx, dst, codegencommentnotice :=
		maybeCtxOpts.newCtx(this, gents), udevgogen.File(this.Name, 2*len(this.Types)*len(gents)), CodeGenCommentNotice.With("progName", CodeGenCommentProgName)
	for _, t := range this.Types {
		for _, g := range gents {
			if ctx.mayGentRunForType(g, t) {
				dst.Body.Add(ctx.generateTopLevelDecls(g, t)...)
			}
		}
	}

	src, timetakengofmt, err = dst.CodeGen(codegencommentnotice, ctx.pkgImportPathsToPkgImportNames, ctx.Opt.EmitNoOpFuncBodies, !ctx.Opt.NoGoFmt)
	timetakengofmt = time.Since(ctx.timeStarted) - timetakengofmt
	return
}

// MustRunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `this`.
func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration) {
	var errs map[*Pkg]error
	timeTakenTotal, timeTakenPerPkg, errs = this.RunGentsAndGenerateOutputFiles(maybeCtxOpts, gents)
	for _, err := range errs {
		panic(err)
	}
	return
}

// RunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `this`.
func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup
	starttime, run := time.Now(), func(pkg *Pkg) {
		src, timetaken, err := pkg.RunGents(maybeCtxOpts, gents)
		if err == nil {
			err = ufs.WriteBinaryFile(filepath.Join(pkg.DirPath, pkg.CodeGen.OutputFileName), src)
		} else {
			println(string(src))
		}
		maps.Lock()
		if timeTakenPerPkg[pkg] = timetaken; err != nil {
			errs[pkg] = err
		}
		maps.Unlock()
		runs.Done()
	}

	timeTakenPerPkg, errs = make(map[*Pkg]time.Duration, len(this)), map[*Pkg]error{}
	runs.Add(len(this))
	for _, pkg := range this {
		go run(pkg)
	}
	runs.Wait()
	timeTakenTotal = time.Since(starttime)
	return
}
