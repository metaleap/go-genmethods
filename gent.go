package gent

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-leap/dev/go/syn"
	"github.com/go-leap/fs"
)

type IGent interface {
	GenerateTopLevelDecls(*Pkg, *Type) []udevgosyn.ISyn
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
	dst := udevgosyn.File(this.Name)
	for _, t := range this.Types {
		for _, g := range gents {
			dst.Body = append(dst.Body, g.GenerateTopLevelDecls(this, t)...)
		}
	}

	emitnoopfuncbodies := EmitNoOpFuncBodies
	if envstr := os.Getenv("GOGENT_EMITNOOPS"); envstr != "" {
		if envbool, e := strconv.ParseBool(envstr); e == nil {
			emitnoopfuncbodies = envbool
		}
	}
	return dst.Src(fmt.Sprintf(CodeGenCommentNotice, CodeGenCommentProgName), emitnoopfuncbodies)
}
