# gent
--
    import "github.com/metaleap/go-gent"

Package gent offers a _Golang code-**gen t**oolkit_; philosophy being:

- _"your package's type-defs. your pluggable custom code-gen logics (+ many
built-in ones), tuned via your struct-field tags. one `go generate` call."_

The design idea is that your codegen programs remains your own `main` packages
written by you, but importing `gent` keeps them short and high-level: fast and
simple to write, iterate, maintain over time. Furthermore (unlike unwieldy
config-file-formats or 100s-of-cmd-args) this approach grants Turing-complete
control over fine-tuning the code-gen flow to only generate what's truly needed,
rather than "every possible func for every possible type-def", to minimize both
codegen and your compilation times.

Focus at the beginning is strictly on generating `func`s and methods for a
package's _existing type-defs_, **not** generating type-defs such as `struct`s.

For building the AST of the to-be-emitted Go source file:

- `gent` relies on my `github.com/go-leap/dev/go/gen` package

- and so do the built-in code-gens under `github.com/metaleap/go-gent/gents/*`,

- but your custom `gent.IGent` implementers are free to prefer other approaches
(such as `text/template` or `github.com/dave/jennifer` or hard-coded
string-building or other) by having their `GenerateTopLevelDecls` implementation
return a `github.com/go-leap/dev/go/gen.SynRaw`-typed byte-array.

Very WIP: more comprehensive readme / package docs to come.

## Usage

```go
var (
	CodeGenCommentNotice   Str = "DON'T EDIT: code gen'd with `{progName}` using `github.com/metaleap/go-gent`"
	CodeGenCommentProgName     = filepath.Base(os.Args[0])

	Defaults struct {
		CtxOpt CtxOpts
	}
	OnBeforeLoad func(*Pkg)
)
```

#### type Ctx

```go
type Ctx struct {
	// options pertaining to this `Ctx`
	Opt CtxOpts

	// strictly read-only
	Pkg *Pkg

	ExtraDefs []*udevgogen.SynFunc

	Gents Gents
}
```

Ctx is a codegen-time context during a `Pkg.RunGents` call and is passed to
`IGent.GenerateTopLevelDecls`.

#### func (*Ctx) DeclsGeneratedSoFar

```go
func (me *Ctx) DeclsGeneratedSoFar(maybeGent IGent, maybeType *Type) (matches []udevgogen.Syns)
```
DeclsGeneratedSoFar collects and returns all results of
`IGent.GenerateTopLevelDecls` performed so far by this `Ctx`, filtered
optionally by `IGent` and/or by `Type`.

#### func (*Ctx) GentExistsFor

```go
func (me *Ctx) GentExistsFor(t *Type, check func(IGent) bool) bool
```

#### func (*Ctx) Import

```go
func (me *Ctx) Import(pkgImportPath string) (pkgImportName udevgogen.PkgName)
```
Import returns the `pkgImportName` for the specified `pkgImportPath`. Eg.
`Import("encoding/json")` might return `pkg__encoding_json` and more
importantly, the import will be properly emitted (only if any of the import's
uses get emitted) at code-gen time. Import is a `Ctx`-local wrapper of the
`github.com/go-leap/dev/go/gen.PkgImports.Ensure` method.

#### func (*Ctx) MayGentRunForType

```go
func (me *Ctx) MayGentRunForType(g IGent, t *Type) bool
```

#### func (*Ctx) N

```go
func (me *Ctx) N(pref string) udevgogen.Named
```

#### type CtxOpts

```go
type CtxOpts struct {
	// For `Defaults.CtxOpts`, initialized from env-var
	// `GOGENT_NOGOFMT` if `strconv.ParseBool`able.
	NoGoFmt bool

	// For `Defaults.CtxOpts`, initialized from env-var
	// `GOGENT_EMITNOOPS` if `strconv.ParseBool`able.
	EmitNoOpFuncBodies bool

	// If set, can be used to prevent running of the given
	// (or any) `IGent` on the given (or any) `*Type`.
	// See also `IGent.Opt().MayRunForType`.
	MayGentRunForType func(IGent, *Type) bool

	// For `Defaults.CtxOpts`, initially set to `"__gent__"`.
	HelpersPrefix string
}
```

CtxOpts wraps `Ctx` options.

#### type Gents

```go
type Gents []IGent
```

Gents is a slice if `IGent`s.

#### func (Gents) EnableOrDisableAll

```go
func (me Gents) EnableOrDisableAll(enabled bool)
```
EnableOrDisableAll sets all `IGent.Opt().Disabled` fields to `!enabled`.

#### func (Gents) EnableOrDisableAllVariantsAndOptionals

```go
func (me Gents) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals calls the same-named method on all
`IGent`s in `me`.

#### func (Gents) With

```go
func (me Gents) With(gents ...Gents) (merged Gents)
```
With merges all `IGent`s in `me` with all those in `gents` into `merged`.

#### type IGent

```go
type IGent interface {
	// must never return `nil` (easiest impl is to embed `Opts`)
	Opt() *Opts

	// implemented as a no-op by `Opts`, to be
	// overridden by implementations as desired
	EnableOrDisableAllVariantsAndOptionals(bool)

	// may read from but never mutate its args.
	// expected to generate preferentially funcs / methods
	// instead of top-level const / var / type decls
	GenerateTopLevelDecls(*Ctx, *Type) udevgogen.Syns
}
```

IGent is the interface implemented by individual code-gens.

#### type Opts

```go
type Opts struct {
	Disabled      bool
	EmitCommented bool

	// not used by `go-gent`, but could be handy for callers
	Name string

	// tested right after `Disabled` and before the below checks.
	// should typically be false for most `Gent`s as the design assumes
	// generation of methods, but it's for the occasional need to generate
	// non-method `func`s related to certain type-alias declarations
	RunOnlyForTypeAliases bool

	RunOnlyOnceWithoutAnyType bool

	// if-and-only-if these are set, they're checked
	// before `MayRunForType` (but after `Disabled`)
	RunNeverForTypes, RunOnlyForTypes struct {
		Named      []string
		Satisfying func(*Ctx, *Type) bool
	}

	// If set, can be used to prevent running of
	// this `IGent` on the given (or any) `*Type`.
	// See also `CtxOpts.MayGentRunForType`.
	MayRunForType func(*Ctx, *Type) bool
}
```

Opts related to a single `IGent`, and designed for embedding.

#### func (*Opts) EnableOrDisableAllVariantsAndOptionals

```go
func (me *Opts) EnableOrDisableAllVariantsAndOptionals(bool)
```
EnableOrDisableAllVariantsAndOptionals implements `IGent` but with a no-op, to
be overridden by `Opts`-embedders as desired.

To disable or enable an `IGent` itself, set `Opts.Disabled`.

#### func (*Opts) Opt

```go
func (me *Opts) Opt() *Opts
```
Opt implements `IGent.Opt()` for `Opts` embedders.

#### type Pkg

```go
type Pkg struct {
	PkgSpec
	DirPath     string
	GoFileNames []string

	Loaded struct {
		Prog *loader.Program
		Info *loader.PackageInfo
	}

	Types Types

	CodeGen struct {
		OutputFile struct {
			Name        string
			DocComments udevgogen.SingleLineDocCommentParagraphs
		}
	}
}
```


#### func  LoadPkg

```go
func LoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string, dontLoadButJustInitUsingPkgNameInstead string) (me *Pkg, err error)
```

#### func  MustLoadPkg

```go
func MustLoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) *Pkg
```

#### func (*Pkg) DirName

```go
func (me *Pkg) DirName() string
```

#### func (*Pkg) RunGents

```go
func (me *Pkg) RunGents(maybeCtxOpts *CtxOpts, gents Gents) (src []byte, stats *Stats, err error)
```
RunGents instructs the given `gents` to generate code for `me`.

#### func (*Pkg) RunGentsAndGenerateOutputFile

```go
func (me *Pkg) RunGentsAndGenerateOutputFile(maybeCtxOpts *CtxOpts, gents Gents) (*Stats, error)
```

#### type PkgSpec

```go
type PkgSpec struct {
	Name       string
	ImportPath string
}
```


#### type Pkgs

```go
type Pkgs map[string]*Pkg
```


#### func  LoadPkgs

```go
func LoadPkgs(pkgPathsWithOutputFileNames map[string]string) (Pkgs, error)
```

#### func  MustLoadPkgs

```go
func MustLoadPkgs(pkgPathsWithOutputFileNames map[string]string) Pkgs
```

#### func (Pkgs) MustRunGentsAndGenerateOutputFiles

```go
func (me Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, statsPerPkg map[*Pkg]*Stats)
```
MustRunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `me`.

#### func (Pkgs) RunGentsAndGenerateOutputFiles

```go
func (me Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpts *CtxOpts, gents Gents) (timeTakenTotal time.Duration, statsPerPkg map[*Pkg]*Stats, errs map[*Pkg]error)
```
RunGentsAndGenerateOutputFiles calls `RunGents` on the `Pkg`s in `me`.

#### type Rename

```go
type Rename func(*Ctx, *Type, string) string
```


#### type Stats

```go
type Stats struct {
	DurationOf struct {
		Constructing time.Duration
		Emitting     time.Duration
		Formatting   time.Duration
		Everything   time.Duration
	}
}
```


#### type Str

```go
type Str string
```


#### func (Str) With

```go
func (me Str) With(placeholderNamesAndValues ...string) string
```

#### type Type

```go
type Type struct {
	Pkg            *Pkg
	SrcFileImports []PkgSpec

	Name  string
	Alias bool

	// Expr is whatever underlying-type this type-decl represents, that is:
	// of the original `type foo bar` or `type foo = bar` declaration,
	// this `Type` is the `foo` identity and its `Expr` captures the `bar`.
	Expr struct {
		// original AST's type-decl's `Expr` (stripped of any&all `ParenExpr`s)
		AstExpr ast.Expr
		// a code-gen `TypeRef` to this `Type` decl's underlying-type
		GenRef *udevgogen.TypeRef
	}

	// commonly useful code-gen values prepared for this `Type`
	G struct {
		// type-ref to this `Type`
		T *udevgogen.TypeRef
		// type-ref to pointer-to-`Type` (think 'ª for addr')
		Tª *udevgogen.TypeRef
		// type-ref to slice-of-`Type`
		Ts *udevgogen.TypeRef
		// type-ref to slice-of-pointers-to-`Type`
		Tªs *udevgogen.TypeRef
		// Name="this" and Type=.G.T
		This udevgogen.NamedTyped
		// Name="this" and Type=.G.Tª
		Thisª udevgogen.NamedTyped
	}

	Enumish struct {
		// expected to be builtin prim-type such as uint8, int64, int --- cases of additional indirections to be handled when they occur in practice
		BaseType string

		ConstNames []string
	}
}
```


#### func (*Type) IsArray

```go
func (me *Type) IsArray() bool
```

#### func (*Type) IsEnumish

```go
func (me *Type) IsEnumish() bool
```

#### func (*Type) IsSlice

```go
func (me *Type) IsSlice() bool
```

#### func (*Type) IsSliceOrArray

```go
func (me *Type) IsSliceOrArray() bool
```

#### func (*Type) SrcFileImportPathByName

```go
func (me *Type) SrcFileImportPathByName(impName string) *PkgSpec
```

#### type Types

```go
type Types []*Type
```


#### func (*Types) Add

```go
func (me *Types) Add(t *Type)
```

#### func (Types) Named

```go
func (me Types) Named(name string) *Type
```

#### type Variant

```go
type Variant struct {
	Add        bool
	DocComment Str
	Name       string
}
```

Variant is like `Variation` but auto-disabled unless `Add` is set.

#### func (*Variant) NameWith

```go
func (me *Variant) NameWith(placeholderNamesAndValues ...string) string
```

#### type Variation

```go
type Variation struct {
	Disabled   bool
	DocComment Str
	Name       string
}
```

Variation is like `Variant` but auto-enabled unless `Disabled` is set.

#### func (*Variation) NameWith

```go
func (me *Variation) NameWith(placeholderNamesAndValues ...string) string
```
