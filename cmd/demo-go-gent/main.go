package main

import (
	"fmt"

	"github.com/metaleap/go-gent"
	// "github.com/metaleap/go-gent/gents/json"
	// "github.com/metaleap/go-gent/gents/maps"
	// "github.com/metaleap/go-gent/gents/trav"
	"github.com/metaleap/go-gent/gents/enum"
)

func init() {
	// gent.Defaults.CtxOpt.EmitNoOpFuncBodies = true
	// gent.Defaults.CtxOpt.NoGoFmt = true
}

func main() {
	pkgs := gent.MustLoadPkgs(map[string]string{
		"github.com/metaleap/go-gent/cmd/demo-go-gent/testpkg": "°_gent.go",
		"github.com/metaleap/zentient":                         "°_gent.go",
		"github.com/metaleap/zentient/lang/golang":             "°_gent.go",
	})

	gents := []gent.IGent{
		&gentenum.Defaults.IsValid,
		// &gentenum.Defaults.IsFoo, // useless & noisy, just a nice simple starting point for custom/new gents
		&gentenum.Defaults.String,
		&gentenum.Defaults.List,
	}

	gent.Defaults.CtxOpt.MayGentRunForType = func(g gent.IGent, t *gent.Type) bool {
		if g == &gentenum.Defaults.String {
			return !(t.Pkg.ImportPath == "github.com/metaleap/zentient" && t.Name == "ToolCats")
		}
		return true
	}

	timetotal, timeperpkg := pkgs.MustRunGentsAndGenerateOutputFiles(nil, gents...)
	fmt.Println("total time taken for all parallel runs and INCL. gofmt + file-write :\n\t\t" + timetotal.String())
	for pkg, timetaken := range timeperpkg {
		fmt.Println("time taken for " + pkg.ImportPath + " EXCL. gofmt & file-write:\n\t\t" + timetaken.String())
	}
}
