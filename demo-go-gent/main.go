package main

import (
	"path/filepath"

	"github.com/go-leap/dev/go/syn"
	"github.com/go-leap/fs"
	"github.com/metaleap/go-gent"
	// "github.com/metaleap/go-gent/gent/json"
	// "github.com/metaleap/go-gent/gent/maps"
	// "github.com/metaleap/go-gent/gent/trav"
	"github.com/metaleap/go-gent/gent/enum"
)

func main() {
	pkgs := map[string]*gent.Pkg{
		"github.com/metaleap/go-gent/demo-go-gent/testpkg": nil,
		"github.com/metaleap/zentient":                     nil,
		"github.com/metaleap/zentient/lang/golang":         nil,
	}
	for pkgpath := range pkgs {
		if pkg, err := gent.LoadPkg(pkgpath, "°gent.go"); err != nil {
			panic(err)
		} else {
			pkgs[pkgpath] = pkg
		}
	}

	gents := []gent.IGent{
		&gentenum.Defaults.Valid,
	}

	udevgosyn.EmitNoOpFuncBodies = false // need occasionally for the odd troubleshoot / pkg-parse-ensuring situation
	for _, pkg := range pkgs {
		src, err := pkg.RunGents(gents...)
		if err == nil {
			err = ufs.WriteBinaryFile(filepath.Join(pkg.DirPath, pkg.OutputFileName), src)
		}
		if err != nil {
			println(string(src))
			panic(err)
		}
	}
}
