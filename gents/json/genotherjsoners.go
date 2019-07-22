package gentjson

import (
	. "github.com/go-leap/dev/go/gen"
	"github.com/metaleap/go-gent"
)

var jsonWriteNull = ˇ.R.Set(B.Append.Of(ˇ.R, "null").Spreads())

func init() {
	Gents.OtherTypes.Marshal.HelpersPrefix, Gents.OtherTypes.Marshal.Name, Gents.OtherTypes.Marshal.DocComment, Gents.OtherTypes.Marshal.InitialBytesCap =
		"jsonMarshal_", DefaultMethodNameMarshal, DefaultDocCommentMarshal, 64
	Gents.OtherTypes.Unmarshal.HelpersPrefix, Gents.OtherTypes.Unmarshal.Name, Gents.OtherTypes.Unmarshal.DocComment =
		"jsonUnmarshal_", DefaultMethodNameUnmarshal, DefaultDocCommentUnmarshal
}

type GentTypeJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
		InitialBytesCap               int
		ResliceInsteadOfWhitespace    bool
		OnStdlibFallbacks             func(*gent.Ctx, ISyn, ...ISyn) Syns
		TryInterfaceTypesBeforeStdlib []*TypeRef
		tryInterfaceTypesDefsDone     bool
	}
	Unmarshal struct {
		JsonMethodOpts
	}
}

func onStdlibDefaultCodegen(_ *gent.Ctx, _ ISyn, s ...ISyn) Syns { return s }

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentTypeJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	if me.Marshal.OnStdlibFallbacks == nil {
		me.Marshal.OnStdlibFallbacks = onStdlibDefaultCodegen
	}
	if !t.IsEnumish() {
		if gennormal, genpanic := me.Marshal.genWhat(t); gennormal || genpanic {
			_ = ctx.N("") // reset counter for dyn-names --- yields less noisy diffs on re-gens
			yield.Add(me.genMarshalMethod(ctx, t, genpanic))
		}
		if gennormal, genpanic := me.Unmarshal.genWhat(t); gennormal || genpanic {
			_ = ctx.N("")
			yield.Add(me.genUnmarshalFromAnyMethod(ctx, t, genpanic),
				me.genUnmarshalMethod(ctx, t, genpanic))
		}
	}
	return
}
