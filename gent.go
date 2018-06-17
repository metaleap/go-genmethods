package gent

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/fs"
)

var MayGentRunForType func(IGent, *Type) bool

type IGent interface {
	GenerateTopLevelDecls(*Type) []udevgogen.ISyn
}

func (this Pkgs) MustRunGentsAndGenerateOutputFiles(gents ...IGent) {
	if err := this.RunGentsAndGenerateOutputFiles(gents...); err != nil {
		panic(err)
	}
}

func (this Pkgs) RunGentsAndGenerateOutputFiles(gents ...IGent) error {
	for _, pkg := range this {
		src, err := pkg.RunGents(gents...)
		if err == nil {
			err = ufs.WriteBinaryFile(filepath.Join(pkg.DirPath, pkg.OutputFileName), src)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *Pkg) RunGents(gents ...IGent) ([]byte, error) {
	dst := udevgogen.File(this.Name)
	for _, t := range this.Types {
		for _, g := range gents {
			if MayGentRunForType == nil || MayGentRunForType(g, t) {
				dst.Body = append(dst.Body, g.GenerateTopLevelDecls(t)...)
			}
		}
	}

	emitnoopfuncbodies := EmitNoOpFuncBodies
	if envstr := os.Getenv("GOGENT_EMITNOOPS"); envstr != "" {
		if envbool, e := strconv.ParseBool(envstr); e == nil {
			emitnoopfuncbodies = envbool
		}
	}
	return dst.Src(fmt.Sprintf(CodeGenCommentNotice, CodeGenCommentProgName), emitnoopfuncbodies, this.CodeGen.PkgImportPathsToPkgImportNames)
}
