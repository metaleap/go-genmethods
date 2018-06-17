package main

import (
	"github.com/metaleap/go-gent"
	// "github.com/metaleap/go-gent/gent/json"
	// "github.com/metaleap/go-gent/gent/maps"
	// "github.com/metaleap/go-gent/gent/trav"
	"github.com/metaleap/go-gent/gent/enum"
)

func init() {
	// gent.EmitNoOpFuncBodies = true
}

func main() {
	pkgs := gent.MustLoadPkgs(map[string]string{
		"github.com/metaleap/go-gent/demo-go-gent/testpkg": "°gent.go",
		"github.com/metaleap/zentient":                     "°gent.go",
		"github.com/metaleap/zentient/lang/golang":         "°gent.go",
	})

	gents := []gent.IGent{
		&gentenum.Defaults.Valid,
		&gentenum.Defaults.IsFoo,
		&gentenum.Defaults.String,
	}

	gent.MayGentRunForType = func(g gent.IGent, t *gent.Type) bool {
		return !(g == &gentenum.Defaults.String && t.Pkg.ImportPath == "github.com/metaleap/zentient" && t.Name == "ToolCats")
	}
	pkgs.MustRunGentsAndGenerateOutputFiles(gents...)
}
