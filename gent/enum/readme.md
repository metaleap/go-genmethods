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
	// Defaults.Valid.MethodName
	MethodName string

	// first of the enumerants
	IsFirstInvalid bool

	// last of the enumerants
	IsLastInvalid bool
}
```

GentValidMethod works for enumish `type`s whose enumerants are ordered such that
the smallest values appear first and the largest last.

#### func (*GentValidMethod) GenerateTopLevelDecls

```go
func (this *GentValidMethod) GenerateTopLevelDecls(pkg *gent.Pkg, t *gent.Type) (tlDecls []gs.IEmit)
```
