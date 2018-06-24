package gent

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/fs"
)

type IGent interface {
	GenerateTopLevelDecls(*Ctx, *Type) udevgogen.Syns
}

func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpt *Opts, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration) {
	var errs map[*Pkg]error
	timeTakenTotal, timeTakenPerPkg, errs = this.RunGentsAndGenerateOutputFiles(maybeCtxOpt, gents...)
	for _, err := range errs {
		panic(err)
	}
	return
}

func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpt *Opts, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup
	starttime, run := time.Now(), func(pkg *Pkg) {
		src, timetaken, err := pkg.RunGents(maybeCtxOpt, gents...)
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

func (this *Pkg) RunGents(maybeCtxOpt *Opts, gents ...IGent) (src []byte, timeTaken time.Duration, err error) {
	dst, codegencommentnotice, ctx :=
		udevgogen.File(this.Name, 2*len(this.Types)*len(gents)), fmt.Sprintf(CodeGenCommentNotice, CodeGenCommentProgName), maybeCtxOpt.newCtx()

	for _, t := range this.Types {
		for _, g := range gents {
			if ctx.Opt.MayGentRunForType == nil || ctx.Opt.MayGentRunForType(g, t) {
				dst.Body.Add(ctx.generateTopLevelDecls(g, t)...)
			}
		}
	}

	src, timeTaken, err = dst.CodeGen(codegencommentnotice, ctx.pkgImportPathsToPkgImportNames, ctx.Opt.EmitNoOpFuncBodies, !ctx.Opt.NoGoFmt)
	timeTaken = time.Since(ctx.timeStarted) - timeTaken
	return
}
