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
	CodeGenCommentNotice   Str = "DO NOT EDIT: code generated with `{progName}` using `github.com/metaleap/go-gent`"
	CodeGenCommentProgName     = filepath.Base(os.Args[0])

	Defaults struct {
		CtxOpt CtxOpts
	}
)
```

#### type Ctx

```go
type Ctx struct {
	Opt CtxOpts
}
```


#### func (*Ctx) DeclsGeneratedSoFar

```go
func (this *Ctx) DeclsGeneratedSoFar(maybeGent IGent, maybeType *Type) (matches []udevgogen.Syns)
```

#### func (*Ctx) I

```go
func (this *Ctx) I(pkgImportPath string) (pkgImportName udevgogen.PkgName)
```

#### type CtxOpts

```go
type CtxOpts struct {
	// For Defaults.CtxOpt, initialized from env-var
	// `GOGENT_NOGOFMT` if `strconv.ParseBool`able.
	NoGoFmt bool

	// For Defaults.CtxOpt, initialized from env-var
	// `GOGENT_EMITNOOPS` if `strconv.ParseBool`able.
	EmitNoOpFuncBodies bool

	// If set, can be used to prevent running of the given
	// (or any) `IGent` on the given (or any) `*Type`.
	MayGentRunForType func(IGent, *Type) bool
}
```


#### type Gents

```go
type Gents []IGent
```


#### func (Gents) EnableOrDisableAll

```go
func (this Gents) EnableOrDisableAll(enabled bool)
```

#### func (Gents) EnableOrDisableAllVariantsAndOptionals

```go
func (this Gents) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```

#### func (Gents) With

```go
func (this Gents) With(gents ...Gents) (all Gents)
```

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
	Disabled              bool
	Name                  string
	RunNeverForTypesNamed []string
	RunOnlyForTypesNamed  []string
	MayRunForType         func(*Type) bool
}
```

Opts related to a single `IGent`, and designed for embedding.

#### func (*Opts) EnableOrDisableAllVariantsAndOptionals

```go
func (this *Opts) EnableOrDisableAllVariantsAndOptionals(bool)
```
EnableOrDisableAllVariantsAndOptionals implements `IGent` but with a no-op, to
be overridden by `Opts`-embedders as desired.

To disable or enable an `IGent` itself, set `Opts.Disabled`.

#### func (*Opts) Opt

```go
func (this *Opts) Opt() *Opts
```
Opt implements `IGent.Opt()` for `Opts` embedders.

#### type Pkg

```go
type Pkg struct {
	Name        string
	ImportPath  string
	DirPath     string
	GoFileNames []string

	Loaded struct {
		Prog *loader.Program
		Info *loader.PackageInfo
	}

	Types Types

	CodeGen struct {
		OutputFileName string
	}
}
```


#### func  LoadPkg

```go
func LoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) (this *Pkg, err error)
```

#### func  MustLoadPkg

```go
func MustLoadPkg(pkgImportPathOrFileSystemPath string, outputFileName string) *Pkg
```

#### func (*Pkg) RunGents

```go
func (this *Pkg) RunGents(maybeCtxOpt *CtxOpts, gents Gents) (src []byte, timeTaken time.Duration, err error)
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
func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpt *CtxOpts, gents Gents) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration)
```

#### func (Pkgs) RunGentsAndGenerateOutputFiles

```go
func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpt *CtxOpts, gents Gents) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error)
```

#### type Str

```go
type Str string
```


#### func (Str) With

```go
func (this Str) With(placeholderNamesAndValues ...string) string
```

#### type Type

```go
type Type struct {
	Pkg *Pkg

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
		// a type-ref to this `Type`
		T *udevgogen.TypeRef
		// a type-ref to pointer-to-`Type`
		Tª *udevgogen.TypeRef
		// a type-ref to slice-of-`Type`
		Ts *udevgogen.TypeRef
		// a type-ref to slice-of-pointers-to-`Type`
		Tªs *udevgogen.TypeRef
		// Name="this" and Type=T.G.T
		This udevgogen.NamedTyped
		// Name="this" and Type=T.G.Tª
		Thisª udevgogen.NamedTyped
	}

	Enumish struct {
		// expected to be builtin prim-type such as uint8, int64, int --- cases of additional indirections to be handled when they occur in practice
		BaseType string

		ConstNames []string
	}
}
```


#### func (*Type) IsEnumish

```go
func (this *Type) IsEnumish() bool
```

#### func (*Type) IsSliceOrArray

```go
func (this *Type) IsSliceOrArray() bool
```

#### type Types

```go
type Types []*Type
```


#### func (*Types) Add

```go
func (this *Types) Add(t *Type)
```

#### func (Types) Named

```go
func (this Types) Named(name string) *Type
```

#### type Variant

```go
type Variant struct {
	Add          bool
	NameOrSuffix string
}
```
