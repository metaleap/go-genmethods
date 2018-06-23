package gent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/fs"
	"github.com/go-leap/sys"
)

var (
	CodeGenCommentNotice   = "DO NOT EDIT: code generated with %s using github.com/metaleap/go-gent"
	CodeGenCommentProgName = filepath.Base(os.Args[0])

	Defaults struct {
		Ctx Ctx
	}

	// If set, can be used to prevent running of the given
	// (or any) `IGent` on the given (or any) `*Type`.
	MayGentRunForType func(IGent, *Type) bool
)

func init() {
	Defaults.Ctx.Opt.EmitNoOpFuncBodies = usys.EnvBool("GOGENT_EMITNOOPS", false)
	Defaults.Ctx.Opt.NoGoFmt = usys.EnvBool("GOGENT_NOGOFMT", false)
}

type IGent interface {
	GenerateTopLevelDecls(*Ctx, *Type) udevgogen.Syns
}

type Str string

func (this Str) With(stringsReplaceOldNew ...string) string {
	return strings.NewReplacer(stringsReplaceOldNew...).Replace(string(this))
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

	declsGenerated map[struct {
		g IGent
		t *Type
	}]udevgogen.Syns
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

func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOptDefaults *Ctx, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration) {
	var errs map[*Pkg]error
	timeTakenTotal, timeTakenPerPkg, errs = this.RunGentsAndGenerateOutputFiles(maybeCtxOptDefaults, gents...)
	for _, err := range errs {
		panic(err)
	}
	return
}

func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOptDefaults *Ctx, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error) {
	var maps sync.Mutex
	var runs sync.WaitGroup
	starttime, run := time.Now(), func(pkg *Pkg) {
		src, timetaken, err := pkg.RunGents(maybeCtxOptDefaults, gents...)
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

func (this *Pkg) RunGents(maybeCtxOptDefaults *Ctx, gents ...IGent) (src []byte, timeTaken time.Duration, err error) {
	dst, codegencommentnotice, ctx :=
		udevgogen.File(this.Name, 2*len(this.Types)*len(gents)), fmt.Sprintf(CodeGenCommentNotice, CodeGenCommentProgName), Ctx{
			Opt: Defaults.Ctx.Opt, TimeStarted: time.Now(), pkgImportPathsToPkgImportNames: udevgogen.PkgImports{},
			declsGenerated: map[struct {
				g IGent
				t *Type
			}]udevgogen.Syns{},
		}
	if maybeCtxOptDefaults != nil {
		ctx.Opt = maybeCtxOptDefaults.Opt
	}

	for _, t := range this.Types {
		for _, g := range gents {
			if MayGentRunForType == nil || MayGentRunForType(g, t) {
				decls := g.GenerateTopLevelDecls(&ctx, t)
				ctx.declsGenerated[struct {
					g IGent
					t *Type
				}{g: g, t: t}] = decls
				dst.Body = append(dst.Body, decls...)
			}
		}
	}

	src, timeTaken, err = dst.CodeGen(codegencommentnotice, ctx.pkgImportPathsToPkgImportNames, ctx.Opt.EmitNoOpFuncBodies, !ctx.Opt.NoGoFmt)
	timeTaken = time.Since(ctx.TimeStarted) - timeTaken
	return
}
