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

GentIsFooMethods generates methods `YourEnumType.IsFoo() bool` for each
enumerant `Foo` in enum type-defs, which equals-compares its receiver to the
respective enumerant `Foo`. (A highly pointless code-gen in real-world terms,
except its exemplary simplicity makes it a handy starter demo sample snippet for
writing new ones from scratch.)

#### func (*GentIsFooMethods) GenerateTopLevelDecls

```go
func (this *GentIsFooMethods) GenerateTopLevelDecls(t *gent.Type) (tlDecls []ISyn)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. If `t` is
a suitable enum type-def, it returns a method `t.IsFoo() bool` for each
enumerant `Foo` in `t`, which equals-compares its receiver to the enumerant.

#### type GentStringMethods

```go
type GentStringMethods struct {
	DocComment string
	Stringers  map[string]func(string) string
	Parsers    struct {
		OnePerStringer         bool
		OneUber                bool
		AddErrlessWithFallback bool
	}
}
```


#### func (*GentStringMethods) GenerateTopLevelDecls

```go
func (this *GentStringMethods) GenerateTopLevelDecls(t *gent.Type) (tlDecls []ISyn)
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
It is only correct for enum type-defs whose enumerants are ordered in the source
such that the numerically smallest values appear first, the largest ones last,
with all enumerant `const`s appearing together.

#### func (*GentValidMethod) GenerateTopLevelDecls

```go
func (this *GentValidMethod) GenerateTopLevelDecls(t *gent.Type) (tlDecls []ISyn)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. It returns
at most one method if `t` is a suitable enum type-def.
