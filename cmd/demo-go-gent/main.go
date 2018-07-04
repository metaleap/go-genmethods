package main

import (
	"fmt"

	"github.com/go-leap/str"
	"github.com/metaleap/go-gent"
	"github.com/metaleap/go-gent/gents/enums"
	"github.com/metaleap/go-gent/gents/slices"
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

	gent.Defaults.CtxOpt.MayGentRunForType = func(g gent.IGent, t *gent.Type) bool {
		if g == &gentenums.Gents.Stringers {
			return !(t.Pkg.ImportPath == "github.com/metaleap/zentient" && t.Name == "ToolCats")
		}
		return true
	}

	gents := gent.Gents{}.With(
		gentenums.Gents.All,
		gentslices.Gents.All,
	)
	gents.EnableOrDisableAllVariantsAndOptionals(true)
	gentenums.Gents.Stringers.All[0].Parser.WithIgnoreCaseCmp = true
	gentenums.Gents.Stringers.All[0].SkipEarlyChecks = true
	gentenums.Gents.Stringers.All[0].EnumerantRename = func(en string) string { return ustr.CaseSnake(ustr.Replace(en, "_", "·")) }

	timetotal, timeperpkg := pkgs.MustRunGentsAndGenerateOutputFiles(nil, gents)
	fmt.Println("total time taken for all parallel runs and INCL. gofmt + file-write :\n\t\t" + timetotal.String())
	for pkg, timetaken := range timeperpkg {
		fmt.Println("time taken for " + pkg.ImportPath + " EXCL. gofmt & file-write:\n\t\t" + timetaken.String())
	}
}
