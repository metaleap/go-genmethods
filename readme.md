# gent
--
    import "github.com/metaleap/go-gent"


## Usage

```go
var (
	CodeGenCommentNotice   = "DO NOT EDIT: code generated with %s using github.com/metaleap/go-gent"
	CodeGenCommentProgName = filepath.Base(os.Args[0])

	Defaults struct {
		Ctx Ctx
	}
)
```

```go
var (
	// If set, can be used to prevent running of the given
	// (or any) `IGent` on the given (or any) `*Type`.
	MayGentRunForType func(IGent, *Type) bool
)
```

#### type Ctx

```go
type Ctx struct {
	// Opt holds the only user-settable fields (in between runs).
	// Code-gens only read but don't mutate them.
	Opt struct {
		NoGoFmt            bool
		EmitNoOpFuncBodies bool
	}

	TimeStarted time.Time
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

#### type IGent

```go
type IGent interface {
	GenerateTopLevelDecls(*Ctx, *Type) udevgogen.Syns
}
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
func (this *Pkg) RunGents(maybeCtxOptDefaults *Ctx, gents ...IGent) (src []byte, timeTaken time.Duration, err error)
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
func (this Pkgs) MustRunGentsAndGenerateOutputFiles(maybeCtxOptDefaults *Ctx, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration)
```

#### func (Pkgs) RunGentsAndGenerateOutputFiles

```go
func (this Pkgs) RunGentsAndGenerateOutputFiles(maybeCtxOptDefaults *Ctx, gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error)
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
