package gent

import (
	"errors"
	"path/filepath"

	"github.com/go-leap/dev/go"
	"github.com/go-leap/fs"
	"github.com/go-leap/str"
	"golang.org/x/tools/go/loader"
)

type Pkg struct {
	OutputFileName string

	Name        string
	ImportPath  string
	DirPath     string
	GoFileNames []string

	Loaded struct {
		Prog *loader.Program
		Info *loader.PackageInfo
	}
}

func LoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) (this *Pkg, err error) {
	errnogopkg := errors.New("not a Go package: " + pkgImportPathOrFileSystemPath)
	this = &Pkg{OutputFileName: outputFileName}

	// figure out DirPath & ImportPath
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

	// figure out *.go files belonging to pkg
	var gofilepaths []string
	if err == nil {
		ufs.WalkFilesIn(this.DirPath, func(fp string) bool {
			if ustr.Suff(fp, ".go") && !ustr.Suff(fp, "_test.go") {
				if fn := filepath.Base(fp); fn != outputFileName {
					gofilepaths, this.GoFileNames = append(gofilepaths, fp), append(this.GoFileNames, fn)
				}
			}
			return true
		})
		if len(this.GoFileNames) == 0 {
			err = errnogopkg
		}
	}

	// actual Go pkg loading
	if err == nil {
		goload := loader.Config{Cwd: this.DirPath}
		goload.CreateFromFilenames(this.ImportPath, gofilepaths...)
		if this.Loaded.Prog, err = goload.Load(); err == nil {
			this.Loaded.Info = this.Loaded.Prog.Package(this.ImportPath)
			for _, gofile := range this.Loaded.Info.Files {
				if gofile.Name != nil {
					if gfname := gofile.Name.Name; this.Name == "" {
						this.Name = gfname
					} else if gfname != "" && gfname != this.Name {
						err = errors.New("naming mismatch: " + this.Name + " vs. " + gfname)
						break
					}
				}
			}
		}
	}

	if err != nil {
		this = nil
	}
	return
}
