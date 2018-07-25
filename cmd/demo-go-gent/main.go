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
	})

	gents := gent.Gents{}.With(
		gentenums.Gents.All,
		gentslices.Gents.All,
	)
	gents.EnableOrDisableAllVariantsAndOptionals(true)
	gentenums.Gents.Stringers.All[0].Parser.WithIgnoreCaseCmp = true
	gentenums.Gents.Stringers.All[0].SkipEarlyChecks = true
	gentenums.Gents.Stringers.All[0].EnumerantRename = func(_ *gent.Ctx, _ *gent.Type, en string) string { return ustr.CaseSnake(ustr.Replace(en, "_", "·")) }

	timetotal, statsperpkg := pkgs.MustRunGentsAndGenerateOutputFiles(nil, gents)
	fmt.Println("total time taken for all parallel runs and INCL. gofmt + file-write :\n\t\t" + timetotal.String())
	for pkg, stats := range statsperpkg {
		fmt.Println("time taken for " + pkg.ImportPath + " EXCL. file-write:\n\t\tconstruct=" + stats.DurationOf.Constructing.String() + "\t\temit=" + stats.DurationOf.Emitting.String() + "\t\tformat=" + stats.DurationOf.Formatting.String() + "")
	}
}
