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
		InternalDecodeMethodName      string
		CommonTypesToExtractToHelpers []*TypeRef
		DefaultCaps                   struct {
			Slices int
			Maps   int
		}
		commonTypesToExtraDefsDone bool
	}

	pkgjson  PkgName
	pkgbytes PkgName
	pkgerrs  PkgName
}

func onStdlibDefaultCodegen(_ *gent.Ctx, _ ISyn, s ...ISyn) Syns { return s }

// GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
func (me *GentTypeJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns) {
	me.pkgjson, me.pkgerrs, me.pkgbytes = ctx.Import("encoding/json"), ctx.Import("errors"), ctx.Import("bytes")
	if me.Marshal.OnStdlibFallbacks == nil {
		me.Marshal.OnStdlibFallbacks = onStdlibDefaultCodegen
	}
	if !t.IsEnumish() {
		if gennormal, genpanic := me.Marshal.genWhat(t); gennormal || genpanic {
			_ = ctx.N("") // reset counter for dyn-names --- yields less noisy diffs on re-gens
			yield.Add(me.genMarshalMethod(ctx, t, genpanic))
		}
		if gennormal, genpanic := me.Unmarshal.genWhat(t); gennormal || genpanic {
			_ = ctx.N("") // see comment above
			yield.Add(me.genUnmarshalMethod(ctx, t, genpanic))
			if gennormal && !genpanic {
				yield.Add(me.genUnmarshalDecodeMethod(ctx, t, genpanic))
			}
		}
	}
	return
}
