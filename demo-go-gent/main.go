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
	// gent.Defaults.Ctx.Opt.EmitNoOpFuncBodies = true
	// gent.Defaults.Ctx.Opt.NoGoFmt = true
}

func main() {
	pkgs := gent.MustLoadPkgs(map[string]string{
		"github.com/metaleap/go-gent/demo-go-gent/testpkg": "°gent.go",
		"github.com/metaleap/zentient":                     "°gent.go",
		"github.com/metaleap/zentient/lang/golang":         "°gent.go",
	})

	gents := []gent.IGent{
		&gentenum.Defaults.IsValid,
		// &gentenum.Defaults.IsFoo, // useless & noisy, disabled by default, just a nice simply starting point for custom/new gents
		&gentenum.Defaults.String,
		&gentenum.Defaults.List,
	}

	gent.MayGentRunForType = func(g gent.IGent, t *gent.Type) bool {
		if g == &gentenum.Defaults.String {
			return !(t.Pkg.ImportPath == "github.com/metaleap/zentient" && t.Name == "ToolCats")
		}
		return true
	}

	timetotal, timeperpkg := pkgs.MustRunGentsAndGenerateOutputFiles(nil, gents...)
	fmt.Println("total time taken for all parallel runs and incl. gofmt + file I/O :\n\t\t" + timetotal.String())
	for pkg, timetaken := range timeperpkg {
		fmt.Println("time taken for " + pkg.ImportPath + " excl. gofmt & file I/O:\n\t\t" + timetaken.String())
	}
}
