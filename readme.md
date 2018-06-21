# gent
--
    import "github.com/metaleap/go-gent"


## Usage

```go
var (
	CodeGenCommentNotice   = "DO NOT EDIT: code generated with %s using github.com/metaleap/go-gent"
	CodeGenCommentProgName = filepath.Base(os.Args[0])

	// overridden by env-var GOGENT_EMITNOOPS, if set to `strconv.ParseBool`able value
	OptEmitNoOpFuncBodies = false

	// overridden by env-var GOGENT_GOFMT, if set to `strconv.ParseBool`able value
	OptGoFmt = true
)
```

```go
var MayGentRunForType func(IGent, *Type) bool
```

#### type IGent

```go
type IGent interface {
	GenerateTopLevelDecls(*Type) udevgogen.Syns
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
		OutputFileName                 string
		PkgImportPathsToPkgImportNames udevgogen.PkgImports
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

#### func (*Pkg) I

```go
func (this *Pkg) I(pkgImportPath string) (pkgImportName string)
```

#### func (*Pkg) RunGents

```go
func (this *Pkg) RunGents(gents ...IGent) (src []byte, timeTaken time.Duration, err error)
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
func (this Pkgs) MustRunGentsAndGenerateOutputFiles(gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration)
```

#### func (Pkgs) RunGentsAndGenerateOutputFiles

```go
func (this Pkgs) RunGentsAndGenerateOutputFiles(gents ...IGent) (timeTakenTotal time.Duration, timeTakenPerPkg map[*Pkg]time.Duration, errs map[*Pkg]error)
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
		BaseType   string
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