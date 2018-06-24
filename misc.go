package gent

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	CodeGenCommentNotice   = "DO NOT EDIT: code generated with `%s` using `github.com/metaleap/go-gent`"
	CodeGenCommentProgName = filepath.Base(os.Args[0])

	Defaults struct {
		CtxOpt Opts
	}
)

func init() {
	CodeGenCommentProgName = strings.TrimPrefix(CodeGenCommentProgName, "zentient-dbg-vsc-go-")
}

type Str string

func (this Str) With(stringsReplaceOldNew ...string) string {
	return strings.NewReplacer(stringsReplaceOldNew...).Replace(string(this))
}
