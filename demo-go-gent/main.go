package main

import (
	"github.com/metaleap/go-gent"
	// "github.com/metaleap/go-gent/gent/json"
	// "github.com/metaleap/go-gent/gent/maps"
	// "github.com/metaleap/go-gent/gent/trav"
)

func main() {
	pkgs := map[string]*gent.Pkg{
		// "github.com/metaleap/go-gent/demo-go-gent/testpkg": nil,
		"github.com/metaleap/zentient": nil,
		// "github.com/metaleap/zentient/lang/golang":         nil,
	}
	for pkgpath := range pkgs {
		if pkg, err := gent.LoadPkg(pkgpath, "°gent.go"); err != nil {
			panic(err)
		} else {
			pkgs[pkgpath] = pkg
			println(len(pkg.Types))
		}
	}
}
