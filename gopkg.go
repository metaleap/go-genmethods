package gent

import (
	"errors"
	"path/filepath"

	"github.com/go-leap/dev/go"
	"github.com/go-leap/dev/go/gen"
	"github.com/go-leap/fs"
	"github.com/go-leap/str"
	"golang.org/x/tools/go/loader"
)

type Pkgs map[string]*Pkg

type PkgSpec struct {
	Name       string
	ImportPath string
}

type Pkg struct {
	PkgSpec
	DirPath     string
	GoFileNames []string

	Loaded struct {
		Prog *loader.Program
		Info *loader.PackageInfo
	}

	Types Types

	CodeGen struct {
		OutputFile struct {
			Name        string
			DocComments udevgogen.SingleLineDocCommentParagraphs
		}
	}
}

func MustLoadPkgs(pkgPathsWithOutputFileNames map[string]string) Pkgs {
	if pkgs, err := LoadPkgs(pkgPathsWithOutputFileNames); err != nil {
		panic(err)
	} else {
		return pkgs
	}
}

func LoadPkgs(pkgPathsWithOutputFileNames map[string]string) (Pkgs, error) {
	pkgs := make(Pkgs, len(pkgPathsWithOutputFileNames))
	for pkgImportPathOrFileSystemPath, outputFileName := range pkgPathsWithOutputFileNames {
		if pkg, err := LoadPkg(pkgImportPathOrFileSystemPath, outputFileName, ""); err != nil {
			return nil, err
		} else {
			pkgs[pkgImportPathOrFileSystemPath] = pkg
		}
	}
	return pkgs, nil
}

func MustLoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) *Pkg {
	if pkg, err := LoadPkg(pkgImportPathOrFileSystemPath, outputFileName, ""); err != nil {
		panic(err)
	} else {
		return pkg
	}
}

func LoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string, dontLoadButJustInitUsingPkgNameInstead string) (me *Pkg, err error) {
	errnogopkg := errors.New("not a Go package: " + pkgImportPathOrFileSystemPath)
	me = &Pkg{PkgSpec: PkgSpec{Name: dontLoadButJustInitUsingPkgNameInstead}}
	me.CodeGen.OutputFile.Name = outputFileName

	if err = me.load_SetPaths(pkgImportPathOrFileSystemPath, errnogopkg); err == nil && dontLoadButJustInitUsingPkgNameInstead == "" {
		if me.CodeGen.OutputFile.Name != "" && OnBeforeLoad != nil {
			OnBeforeLoad(me)
		}
		var gofilepaths []string
		if gofilepaths, err = me.load_SetFileNames(errnogopkg); err == nil {
			err = me.load_FromFiles(gofilepaths)
		}
	}
	if err != nil {
		me = nil
	}
	return
}

func (me *Pkg) load_SetPaths(pkgImportPathOrFileSystemPath string, errnogopkg error) (err error) {
	if ufs.IsDir(pkgImportPathOrFileSystemPath) {
		me.DirPath = pkgImportPathOrFileSystemPath
	} else if ufs.IsFile(pkgImportPathOrFileSystemPath) {
		me.DirPath = filepath.Dir(pkgImportPathOrFileSystemPath)
	} else if me.DirPath = udevgo.GopathSrc(pkgImportPathOrFileSystemPath); me.DirPath != "" && ufs.IsDir(me.DirPath) {
		me.ImportPath = pkgImportPathOrFileSystemPath
	} else {
		err = errnogopkg
	}
	if err == nil && me.ImportPath == "" && me.DirPath != "" {
		me.ImportPath = udevgo.DirPathToImportPath(me.DirPath)
	}
	if err == nil && me.ImportPath == "" {
		err = errnogopkg
	}
	return
}

func (me *Pkg) load_SetFileNames(errnogopkg error) (goFilePaths []string, err error) {
	ufs.WalkFilesIn(me.DirPath, func(fp string) bool {
		if fn := filepath.Base(fp); ustr.Suff(fp, ".go") && !ustr.Suff(fp, "_test.go") {
			goFilePaths, me.GoFileNames = append(goFilePaths, fp), append(me.GoFileNames, fn)
		}
		return true
	})
	if len(me.GoFileNames) == 0 {
		err = errnogopkg
	}
	return
}

func (me *Pkg) load_FromFiles(goFilePaths []string) (err error) {
	goload := loader.Config{Cwd: me.DirPath}
	goload.CreateFromFilenames(me.ImportPath, goFilePaths...)
	if me.Loaded.Prog, err = goload.Load(); err == nil {
		me.Loaded.Info = me.Loaded.Prog.Package(me.ImportPath)
		for _, gofile := range me.Loaded.Info.Files {
			if gofile.Name != nil {
				if gfname := gofile.Name.Name; me.Name == "" {
					me.Name = gfname
				} else if gfname != "" && gfname != me.Name {
					err = errors.New("naming mismatch: " + me.Name + " vs. " + gfname)
					return
				}
				me.load_Types(gofile)
			}
		}
		me.load_PopulateTypes()
	}
	return
}

func (me *Pkg) DirName() string { return filepath.Base(me.DirPath) }
