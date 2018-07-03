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

type Gents []IGent

func (this Gents) With(gents ...Gents) (all Gents) {
	all = append(make(Gents, 0, len(this)+2*len(gents)), this...)
	for i := range gents {
		all = append(all, gents[i]...)
	}
	return
}

func (this Gents) EnableOrDisableAll(enabled bool) {
	disabled := !enabled
	for i := range this {
		this[i].Opt().Disabled = disabled
	}
}

func (this Gents) EnableOrDisableAllVariantsAndOptionals(enabled bool) {
	for i := range this {
		this[i].EnableOrDisableAllVariantsAndOptionals(enabled)
	}
}

// Opts related to a single `IGent`, and designed for embedding.
type Opts struct {
	Disabled              bool
	Name                  string
	RunNeverForTypesNamed []string
	RunOnlyForTypesNamed  []string
	MayRunForType         func(*Type) bool
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

func (this *Pkg) RunGents(maybeCtxOpt *CtxOpts, gents Gents) (src []byte, timeTaken time.Duration, err error) {
	dst, ctx, codegencommentnotice :=
		udevgogen.File(this.Name, 2*len(this.Types)*len(gents)), maybeCtxOpt.newCtx(this, gents), CodeGenCommentNotice.With("progName", CodeGenCommentProgName)
	for _, t := range this.Types {
		for _, g := range gents {
			if ctx.mayGentRunForType(g, t) {
				dst.Body.Add(ctx.generateTopLevelDecls(g, t)...)
			}
		}
	}

	src, timeTaken, err = dst.CodeGen(codegencommentnotice, ctx.pkgImportPathsToPkgImportNames, ctx.Opt.EmitNoOpFuncBodies, !ctx.Opt.NoGoFmt)
	timeTaken = time.Since(ctx.timeStarted) - timeTaken
	return
}

func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpt *CtxOpts, gents Gents) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration) {
	var errs map[*Pkg]error
	timeTakenTotal, timeTakenPerPkg, errs = this.RunGentsAndGenerateOutputFiles(maybeCtxOpt, gents)
	for _, err := range errs {
		panic(err)
	}
	return
}

func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpt *CtxOpts, gents Gents) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup
	starttime, run := time.Now(), func(pkg *Pkg) {
		src, timetaken, err := pkg.RunGents(maybeCtxOpt, gents)
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
