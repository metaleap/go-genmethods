package main

import (
	"github.com/metaleap/go-gent"
	// "github.com/metaleap/go-gent/gent/json"
)

func main() {
	pkgpath := "github.com/metaleap/go-gent/demo-go-gent/testpkg"
	pkg, err := gent.LoadPkg(pkgpath, "°gent.go")
	if err != nil {
		panic(err)
	}
	println(len(pkg.Types))
}
