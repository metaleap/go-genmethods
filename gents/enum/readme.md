# gentenum
--
    import "github.com/metaleap/go-gent/gents/enum"

Package gentenum provides `gent.IGent` code-gens of `func`s related to "enum-ish
type-defs". Most of them expect and assume enum type-defs whose enumerants are
ordered in the source such that the numerically smallest values appear first,
the largest ones last, with all enumerant `const`s appearing next to each other.

## Usage

```go
var (
	Defaults struct {
		IsValid GentIsValidMethod
		IsFoo   GentIsFooMethods
		String  GentStringMethods
		List    GentListEnumerantsFunc
	}
)
```

#### type GentIsFooMethods

```go
type GentIsFooMethods struct {
	Disabled   bool
	DocComment gent.Str

	// eg `Is{e}` -> `IsMyOne`, `IsMyTwo`, etc.
	MethodName gent.Str

	// if set, renames the enumerant used for {e} in `MethodName`
	MethodNameRenameEnumerant func(string) string
}
```

GentIsFooMethods generates methods `YourEnumType.IsFoo() bool` for each
enumerant `Foo` in enum type-defs, which equals-compares its receiver to the
respective enumerant `Foo`. (A HIGHLY POINTLESS code-gen in real-world terms,
except its exemplary simplicity makes it a handy
starter-demo-sample-snippet-blueprint for writing new ones from scratch.)

An instance with illustrative defaults is in `Defaults.IsFoo`.

#### func (*GentIsFooMethods) GenerateTopLevelDecls

```go
func (this *GentIsFooMethods) GenerateTopLevelDecls(t *gent.Type) (decls Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. If `t` is
a suitable enum type-def, it returns a method `t.IsFoo() bool` for each
enumerant `Foo` in `t`, which equals-compares its receiver to the enumerant.

#### type GentIsValidMethod

```go
type GentIsValidMethod struct {
	Disabled       bool
	DocComment     gent.Str
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}
```

GentIsValidMethod generates a `Valid` method for enum type-defs, which checks
whether the receiver value seems to be within the range of the known enumerants.

An instance with illustrative defaults is in `Defaults.IsValid`.

#### func (*GentIsValidMethod) GenerateTopLevelDecls

```go
func (this *GentIsValidMethod) GenerateTopLevelDecls(t *gent.Type) (decls Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. It returns
at most one method if `t` is a suitable enum type-def.

#### type GentListEnumerantsFunc

```go
type GentListEnumerantsFunc struct {
	Disabled   bool
	DocComment gent.Str

	// eg. "Wellknown{T}{s}" with `{T}` for type name and
	// `{s}` for pluralization suffix (either "s" or "es")
	FuncName gent.Str
}
```

GentListEnumerantsFunc generates a `func WellknownFoos() ([]string, []Foo)` for
each enum type-def `Foo`.

An instance with illustrative defaults is in `Defaults.List`.

#### func (*GentListEnumerantsFunc) GenerateTopLevelDecls

```go
func (this *GentListEnumerantsFunc) GenerateTopLevelDecls(t *gent.Type) (decls Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentStringMethods

```go
type GentStringMethods struct {
	Disabled    bool
	Stringers   []Stringer
	DocComments struct {
		Parsers               gent.Str
		ParsersErrlessVariant gent.Str
	}
}
```

GentStringMethods generates for enum type-defs the specified `string`ifying
methods, optionally with corresponding "parsing" funcs.

An instance with illustrative defaults is in `Defaults.String`.

#### func (*GentStringMethods) GenerateTopLevelDecls

```go
func (this *GentStringMethods) GenerateTopLevelDecls(t *gent.Type) (decls Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type Stringer

```go
type Stringer struct {
	Disabled                         bool
	DocComment                       gent.Str
	Name                             string
	EnumerantRename                  func(string) string
	ParseFuncName                    gent.Str
	ParseAddIgnoreCaseCmp            bool
	ParseAddErrlessVariantWithSuffix string
}
```
