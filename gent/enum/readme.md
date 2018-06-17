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
	DocComment       string
	MethodNamePrefix string
	RenameEnumerant  func(string) string
}
```

GentIsFooMethods generates a method `t.IsFoo() bool` for each enumerant `Foo` in
enums, which equals-compares its receiver. (A hugely pointless code-gen, but its
simplicity makes it a decent starter example for writing custom ones.)

#### func (*GentIsFooMethods) GenerateTopLevelDecls

```go
func (this *GentIsFooMethods) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. If `t` is
a suitable enum type-def, it returns a method `t.IsFoo() bool` for each
enumerant `Foo` in `t`, which equals-compares its receiver.

#### type GentStringMethods

```go
type GentStringMethods struct {
	DocComment string
}
```


#### type GentValidMethod

```go
type GentValidMethod struct {
	DocComment     string
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}
```

GentValidMethod generates a `Valid` method for enum type-defs, which checks
whether the receiver value seems to be within the range of the known enumerants.
It's only correct for enum type-defs whose enumerants are ordered in the source
such that the smallest values appear first, the largest last, and with all
enumerant `const`s appearing together.

#### func (*GentValidMethod) GenerateTopLevelDecls

```go
func (this *GentValidMethod) GenerateTopLevelDecls(_ *gent.Pkg, t *gent.Type) (tlDecls []ISyn)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. It returns
at most one method if `t` is a suitable enum type-def.
