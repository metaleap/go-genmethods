package gentenums

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

// GentIsValidMethod generates a `Valid` method for enum type-defs, which checks
// whether the receiver value seems to be within the range of the known enumerants.
//
// An instance with illustrative defaults is in `Defaults.IsValid`.
type GentIsValidMethod struct {
	gent.Opts

	DocComment     gent.Str
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
// It returns at most one method if `t` is a suitable enum type-def.
func (this *GentIsValidMethod) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns) {
	if t.IsEnumish() {
		firstinvalid, firstname, lastname, firsthint, lasthint :=
			this.IsFirstInvalid, t.Enumish.ConstNames[0], t.Enumish.ConstNames[len(t.Enumish.ConstNames)-1], "inclusive", "inclusive"
		if firstname == "_" {
			firstinvalid, firstname = false, t.Enumish.ConstNames[1]
		}

		var firstoperator, lastoperator ISyn = Geq(V.This, N(firstname)), Leq(V.This, N(lastname))
		if firstinvalid {
			firstoperator, firsthint = Gt(V.This, N(firstname)), "exclusive"
		}
		if this.IsLastInvalid {
			lastoperator, lasthint = Lt(V.This, N(lastname)), "exclusive"
		}

		method := Fn(t.CodeGen.ThisVal, this.MethodName, &Sigs.NoneToBool,
			Set(V.Ret, And(firstoperator, lastoperator)),
		)
		method.Doc.Add(this.DocComment.With(
			"{N}", this.MethodName,
			"{T}", t.Name,
			"{fn}", firstname,
			"{fh}", firsthint,
			"{ln}", lastname,
			"{lh}", lasthint,
		))
		decls = Syns{method}
	}
	return
}
