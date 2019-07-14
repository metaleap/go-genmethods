package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

var jsonWriteNull = ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads())

func init() {
	Gents.OtherTypes.Marshal.Name, Gents.OtherTypes.Marshal.DocComment, Gents.OtherTypes.Marshal.InitialBytesCap =
		DefaultMethodNameMarshal, DefaultDocCommentMarshal, 64
	Gents.OtherTypes.Unmarshal.Name, Gents.OtherTypes.Unmarshal.DocComment =
		DefaultMethodNameUnmarshal, DefaultDocCommentUnmarshal
}

type GentTypeJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
		InitialBytesCap               int
		ResliceInsteadOfWhitespace    bool
		GenPrintlnOnStdlibFallbacks   bool
		TryInterfaceTypesBeforeStdlib []*TypeRef
		tryInterfaceTypesDefsDone     bool
	}
	Unmarshal struct {
		JsonMethodOpts
	}
}

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentTypeJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if !t.IsEnumish() {
		if gmn, gmp := me.Marshal.genWhat(t); gmn || gmp {
			yield.Add(me.genMarshalMethod(ctx, t, gmp))
		}
		if gun, gup := me.Unmarshal.genWhat(t); gun || gup {
			yield.Add(me.genUnmarshalMethod(ctx, t, gup))
		}
	}
	return
}
