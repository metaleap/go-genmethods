package gent

import (
	"errors"
	"path/filepath"

	"github.com/go-leap/dev/go"
	"github.com/go-leap/fs"
	"github.com/go-leap/str"
	"golang.org/x/tools/go/loader"
)

type Pkgs map[string]*Pkg

type Pkg struct {
	Name        string
	ImportPath  string
	DirPath     string
	GoFileNames []string

	Loaded struct {
		Prog *loader.Program
		Info *loader.PackageInfo
	}

	Types Types

	CodeGen struct {
		OutputFileName string
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
		if pkg, err := LoadPkg(pkgImportPathOrFileSystemPath, outputFileName); err != nil {
			return nil, err
		} else {
			pkgs[pkgImportPathOrFileSystemPath] = pkg
		}
	}
	return pkgs, nil
}

func MustLoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) *Pkg {
	if pkg, err := LoadPkg(pkgImportPathOrFileSystemPath, outputFileName); err != nil {
		panic(err)
	} else {
		return pkg
	}
}

func LoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) (this *Pkg, err error) {
	errnogopkg := errors.New("not a Go package: " + pkgImportPathOrFileSystemPath)
	this = &Pkg{}
	this.CodeGen.OutputFileName = outputFileName

	if err = this.load_SetPaths(pkgImportPathOrFileSystemPath, errnogopkg); err == nil {
		var gofilepaths []string
		if gofilepaths, err = this.load_SetFileNames(errnogopkg); err == nil {
			err = this.load_FromFiles(gofilepaths)
		}
	}

	if err != nil {
		this = nil
	}
	return
}

func (this *Pkg) load_SetPaths(pkgImportPathOrFileSystemPath string, errnogopkg error) (err error) {
	if ufs.IsDir(pkgImportPathOrFileSystemPath) {
		this.DirPath = pkgImportPathOrFileSystemPath
	} else if ufs.IsFile(pkgImportPathOrFileSystemPath) {
		this.DirPath = filepath.Dir(pkgImportPathOrFileSystemPath)
	} else if this.DirPath = udevgo.GopathSrc(pkgImportPathOrFileSystemPath); this.DirPath != "" && ufs.IsDir(this.DirPath) {
		this.ImportPath = pkgImportPathOrFileSystemPath
	} else {
		err = errnogopkg
	}
	if err == nil && this.ImportPath == "" && this.DirPath != "" {
		this.ImportPath = udevgo.DirPathToImportPath(this.DirPath)
	}
	if err == nil && this.ImportPath == "" {
		err = errnogopkg
	}
	return
}

func (this *Pkg) load_SetFileNames(errnogopkg error) (goFilePaths []string, err error) {
	ufs.WalkFilesIn(this.DirPath, func(fp string) bool {
		if fn := filepath.Base(fp); ustr.Suff(fp, ".go") && !ustr.Suff(fp, "_test.go") {
			goFilePaths, this.GoFileNames = append(goFilePaths, fp), append(this.GoFileNames, fn)
		}
		return true
	})
	if len(this.GoFileNames) == 0 {
		err = errnogopkg
	}
	return
}

func (this *Pkg) load_FromFiles(goFilePaths []string) (err error) {
	goload := loader.Config{Cwd: this.DirPath}
	goload.CreateFromFilenames(this.ImportPath, goFilePaths...)
	if this.Loaded.Prog, err = goload.Load(); err == nil {
		this.Loaded.Info = this.Loaded.Prog.Package(this.ImportPath)
		for _, gofile := range this.Loaded.Info.Files {
			if gofile.Name != nil {
				if gfname := gofile.Name.Name; this.Name == "" {
					this.Name = gfname
				} else if gfname != "" && gfname != this.Name {
					err = errors.New("naming mismatch: " + this.Name + " vs. " + gfname)
					return
				}
				this.load_Types(gofile)
			}
		}
	}
	return
}
