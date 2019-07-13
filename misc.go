package gent

import (
	"os"
	"path/filepath"

	"github.com/go-leap/str"
)

var (
	CodeGenCommentNotice   Str = "DON'T EDIT: code gen'd with `{progName}` using `github.com/metaleap/go-gent`"
	CodeGenCommentProgName     = filepath.Base(os.Args[0])

	Defaults struct {
		CtxOpt CtxOpts
	}
	OnBeforeLoad func(*Pkg)
)

func init() {
	CodeGenCommentProgName = ustr.TrimPref(CodeGenCommentProgName, "zentient-dbg-vsc-go-")
}

type Str string

func (me Str) With(placeholderNamesAndValues ...string) string {
	return strWith(string(me), placeholderNamesAndValues...)
}

var strWith = ustr.NamedPlaceholders('{', '}')

// Variant is like `Variation` but auto-disabled unless `Add` is set.
type Variant struct {
	Add        bool
	DocComment Str
	Name       string
}

func (me *Variant) NameWith(placeholderNamesAndValues ...string) string {
	return Str(me.Name).With(placeholderNamesAndValues...)
}

// Variation is like `Variant` but auto-enabled unless `Disabled` is set.
type Variation struct {
	Disabled   bool
	DocComment Str
	Name       string
}

func (me *Variation) NameWith(placeholderNamesAndValues ...string) string {
	return Str(me.Name).With(placeholderNamesAndValues...)
}

type Rename func(*Ctx, *Type, string) string
