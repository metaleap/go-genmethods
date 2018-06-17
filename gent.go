package gent

import (
	"github.com/go-leap/dev/go/syn"
)

type IGent interface {
	GenerateTopLevelDecls(*Pkg, *Type) []udevgosyn.IEmit
}
