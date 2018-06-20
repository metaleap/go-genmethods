package gent

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/fs"
	"github.com/go-leap/sys"
)

var MayGentRunForType func(IGent, *Type) bool

type IGent interface {
	GenerateTopLevelDecls(*Type) udevgogen.Syns
}

func (this Pkgs) MustRunGentsAndGenerateOutputFiles(gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration) {
	var errs map[*Pkg]error
	timeTakenTotal, timeTakenPerPkg, errs = this.RunGentsAndGenerateOutputFiles(gents...)
	for _, err := range errs {
		panic(err)
	}
	return
}

func (this Pkgs) RunGentsAndGenerateOutputFiles(gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup
	starttime, run := time.Now(), func(pkg *Pkg) {
		src, timetaken, err := pkg.RunGents(gents...)
		if err == nil {
			err = ufs.WriteBinaryFile(filepath.Join(pkg.DirPath, pkg.CodeGen.OutputFileName), src)
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

func (this *Pkg) RunGents(gents ...IGent) (src []byte, timeTaken time.Duration, err error) {
	dst, codegencommentnotice := udevgogen.File(this.Name, 2*len(this.Types)*len(gents)), fmt.Sprintf(CodeGenCommentNotice, CodeGenCommentProgName)
	optnoops, optgofmt := usys.EnvBool("GOGENT_EMITNOOPS", OptEmitNoOpFuncBodies), usys.EnvBool("GOGENT_GOFMT", OptGoFmt)

	timestarted := time.Now()
	for _, t := range this.Types {
		for _, g := range gents {
			if MayGentRunForType == nil || MayGentRunForType(g, t) {
				dst.Body = append(dst.Body, g.GenerateTopLevelDecls(t)...)
			}
		}
	}

	src, timeTaken, err = dst.CodeGen(codegencommentnotice, this.CodeGen.PkgImportPathsToPkgImportNames, optnoops, optgofmt)
	timeTaken = time.Since(timestarted) - timeTaken
	return
}
