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
	CodeGenCommentNotice   = "DO NOT EDIT: code generated with `%s` using `github.com/metaleap/go-gent`"
	CodeGenCommentProgName = filepath.Base(os.Args[0])

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
func (this *Ctx) I(pkgImportPath string) (pkgImportName string)
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


#### type IGent

```go
type IGent interface {
	Opt() *Opts
	GenerateTopLevelDecls(*Ctx, *Type) udevgogen.Syns
}
```


#### type Opts

```go
type Opts struct {
	Disabled bool
}
```


#### func (*Opts) Opt

```go
func (this *Opts) Opt() *Opts
```

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
func (this *Pkg) RunGents(maybeCtxOpt *CtxOpts, gents ...IGent) (src []byte, timeTaken time.Duration, err error)
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
func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOpt *CtxOpts, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration)
```

#### func (Pkgs) RunGentsAndGenerateOutputFiles

```go
func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOpt *CtxOpts, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error)
```

#### type Str

```go
type Str string
```


#### func (Str) With

```go
func (this Str) With(stringsReplaceOldNew ...string) string
```

#### type Type

```go
type Type struct {
	Pkg *Pkg

	Name  string
	Decl  *ast.TypeSpec
	Alias bool

	Ast struct {
		Named      *ast.Ident
		Imported   *ast.SelectorExpr
		Ptr        *ast.StarExpr
		TArrOrSl   *ast.ArrayType
		TChan      *ast.ChanType
		TFunc      *ast.FuncType
		TInterface *ast.InterfaceType
		TMap       *ast.MapType
		TStruct    *ast.StructType
	}

	Enumish struct {
		// expected to be builtin prim-type such as uint8, int64, int --- cases of additional indirections to be handled when they occur in practice
		BaseType string

		ConstNames []string
	}

	CodeGen struct {
		ThisVal udevgogen.NamedTyped
		ThisPtr udevgogen.NamedTyped
		Ref     *udevgogen.TypeRef
	}
}
```


#### func (*Type) SeemsEnumish

```go
func (this *Type) SeemsEnumish() bool
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

#### func (Types) Struct

```go
func (this Types) Struct(name string) *Type
```
