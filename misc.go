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

type Variant struct {
	Add        bool
	Name       string
	DocComment Str
}

func (this *Variant) NameWith(placeholderNamesAndValues ...string) string {
	return Str(this.Name).With(placeholderNamesAndValues...)
}

type Variation struct {
	Disabled   bool
	DocComment Str
	Name       string
}

func (this *Variation) NameWith(placeholderNamesAndValues ...string) string {
	return Str(this.Name).With(placeholderNamesAndValues...)
}

type Rename func(*Ctx, *Type, string) string
