package gent

import (
	"os"
	"path/filepath"

	"github.com/go-leap/str"
)

var (
	CodeGenCommentNotice   Str = "DO NOT EDIT: code generated with `{progName}` using `github.com/metaleap/go-gent`"
	CodeGenCommentProgName     = filepath.Base(os.Args[0])

	Defaults struct {
		CtxOpt CtxOpts
	}
)

func init() {
	CodeGenCommentProgName = ustr.TrimPref(CodeGenCommentProgName, "zentient-dbg-vsc-go-")
}

type Str string

func (this Str) With(placeholderNamesAndValues ...string) string {
	return strWith(string(this), placeholderNamesAndValues...)
}

var strWith = ustr.NamedPlaceholders('{', '}')

type Variadic bool

type Variant struct {
	Add          bool
	NameOrSuffix string
}
