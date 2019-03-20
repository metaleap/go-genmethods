# gentenums
--
    import "github.com/metaleap/go-gent/gents/enums"

Package gentenums provides `gent.IGent` code-gens of `func`s related to
"enum-ish type-defs". Most of them expect and assume enum type-defs whose
enumerants are ordered in the source such that the numerically smallest value
appears first, the largest one last, with all enumerant `const`s appearing next
to each other.

## Usage

```go
const (
	DefaultIsFooDocComment = "{N} returns whether the value of this `{T}` equals `{e}`."
	DefaultIsFooMethodName = "Is{e}"
)
```

```go
const (
	DefaultListBothFuncName     = "Wellknown{T}{s}"
	DefaultListBothDocComment   = "{N} returns the `names` and `values` of all {n} well-known `{T}` enumerants."
	DefaultListNamesFuncName    = "Wellknown{T}Names"
	DefaultListNamesDocComment  = "{N} returns the `names` of all {n} well-known `{T}` enumerants."
	DefaultListValuesFuncName   = "Wellknown{T}Values"
	DefaultListValuesDocComment = "{N} returns the `values` of all {n} well-known `{T}` enumerants."
	DefaultListMapFuncName      = "Wellknown{T}NamesAndValues"
	DefaultListMapDocComment    = "{N} returns the `namesToValues` of all {n} well-known `{T}` enumerants."
)
```

```go
const (
	DefaultStringers0DocComments              = "{N} implements the Go standard library's `fmt.Stringer` interface."
	DefaultStringers0MethodName               = "String"
	DefaultStringers1DocComments              = "{N} implements the Go standard library's `fmt.GoStringer` interface."
	DefaultStringers1MethodName               = "GoString"
	DefaultStringersParsersDocComments        = "{N} returns the `{T}` represented by `{s}` (as returned by `{T}.{str}`, {caseSensitivity}), or an `error` if none exists."
	DefaultStringersParsersDocCommentsErrless = "{N} is like `{p}` but returns `{fallback}` for unrecognized inputs."
	DefaultStringersParsersFuncName           = "{T}From{str}"
)
```

```go
const (
	DefaultIsValidDocComment = "{N} returns whether the value of this `{T}` is between `{fn}` ({fh}) and `{ln}` ({lh})."
	DefaultIsValidMethodName = "Valid"
)
```

```go
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IsFoo     GentIsFooMethods
		IsValid   GentIsValidMethod
		Listers   GentListEnumerantsFuncs
		Stringers GentStringersMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

#### type GentIsFooMethods

```go
type GentIsFooMethods struct {
	gent.Opts

	DocComment gent.Str
	// eg `Is{e}` -> `IsMyOne`, `IsMyTwo`, etc.
	MethodName gent.Str

	// if set, renames the enumerant used for {e} in `MethodName`
	MethodNameRenameEnumerant gent.Rename
}
```

GentIsFooMethods generates methods `YourEnumType.IsFoo() bool` for each
enumerant `Foo` in enum type-defs, which equals-compares its receiver to the
respective enumerant `Foo`. (A HIGHLY POINTLESS code-gen in real-world terms,
except its exemplary simplicity makes it a handy
starter-demo-sample-snippet-blueprint for writing new ones from scratch.)

An instance with illustrative defaults is in `Gents.IsFoo`.

#### func (*GentIsFooMethods) GenerateTopLevelDecls

```go
func (me *GentIsFooMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. If `t` is
a suitable enum type-def, it returns a method `t.IsFoo() bool` for each
enumerant `Foo` in `t`, which equals-compares its receiver to the enumerant.

#### type GentIsValidMethod

```go
type GentIsValidMethod struct {
	gent.Opts

	DocComment     gent.Str
	MethodName     string
	IsFirstInvalid bool
	IsLastInvalid  bool
}
```

GentIsValidMethod generates a `Valid` method for enum type-defs, which checks
whether the receiver value seems to be within the range of the known enumerants.

An instance with illustrative defaults is in `Gents.IsValid`.

#### func (*GentIsValidMethod) GenerateTopLevelDecls

```go
func (me *GentIsValidMethod) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`. It returns
at most one method if `t` is a suitable enum type-def.

#### type GentListEnumerantsFuncs

```go
type GentListEnumerantsFuncs struct {
	gent.Opts

	ListNames  gent.Variation
	ListValues gent.Variation
	ListMap    gent.Variation
	ListBoth   gent.Variation

	AlwaysSkipFirst bool
	Rename          gent.Rename
}
```

GentListEnumerantsFuncs generates a `func WellknownFoos() ([]string, []Foo)` for
each enum type-def `Foo`.

An instance with illustrative defaults is in `Gents.List`.

#### func (*GentListEnumerantsFuncs) EnableOrDisableAllVariantsAndOptionals

```go
func (me *GentListEnumerantsFuncs) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals implements
`github.com/metaleap/go-gent.IGent`.

#### func (*GentListEnumerantsFuncs) GenerateTopLevelDecls

```go
func (me *GentListEnumerantsFuncs) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentStringersMethods

```go
type GentStringersMethods struct {
	gent.Opts

	All         []StringMethodOpts
	DocComments struct {
		Parsers               gent.Str
		ParsersErrlessVariant gent.Str
	}
}
```

GentStringersMethods generates for enum type-defs the specified `string`ifying
methods, optionally with corresponding "parsing" funcs.

An instance with illustrative defaults is in `Gents.String`.

#### func (*GentStringersMethods) EnableOrDisableAllVariantsAndOptionals

```go
func (me *GentStringersMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals implements
`github.com/metaleap/go-gent.IGent`.

#### func (*GentStringersMethods) GenerateTopLevelDecls

```go
func (me *GentStringersMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type StringMethodOpts

```go
type StringMethodOpts struct {
	gent.Variation
	EnumerantRename gent.Rename
	SkipEarlyChecks bool
	Parser          struct {
		Add               bool
		WithIgnoreCaseCmp bool
		Errless           gent.Variant
		FuncName          gent.Str
	}
}
```
