# gentenum
--
    import "github.com/metaleap/go-gent/gent/enum"


## Usage

```go
var (
	Defaults struct {
		Valid  GentValidMethod
		IsFoo  GentIsFooMethods
		String GentStringMethods
	}
)
```

#### type GentIsFooMethods

```go
type GentIsFooMethods struct {
}
```


#### type GentStringMethods

```go
type GentStringMethods struct {
}
```


#### type GentValidMethod

```go
type GentValidMethod struct {
	// defaults to Defaults.Valid.DocComment
	DocComment string

	// defaults to Defaults.Valid.MethodName
	MethodName string

	// if `true`, generate gt instead of geq
	IsFirstInvalid bool

	// if `true`, generate lt instead of leq
	IsLastInvalid bool
}
```

GentValidMethod generated a `Valid` method for enum type-defs, checking whether
the value seems to be within the range of the known enumerants. It only supports
enum type-defs whose enumerants are ordered in the source such that the smallest
values appear first and the largest last, with all enumerant `const`s appearing
together.

#### func (*GentValidMethod) GenerateTopLevelDecls

```go
func (this *GentValidMethod) GenerateTopLevelDecls(pkg *gent.Pkg, t *gent.Type) (tlDecls []gs.IEmit)
```
GenerateTopLevelDecls implements github.com/metaleap/go-gent.IGent
