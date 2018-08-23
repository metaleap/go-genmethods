# gentslices
--
    import "github.com/metaleap/go-gent/gents/slices"


## Usage

```go
const (
	DefaultConvFieldsDocComment = "{N} returns all `{field}` values of the constituent `{T}`s in `{this}`."
	DefaultConvToMapsDocComment = "{N} converts `{this}` into a `map` indexed by the `{field}` values of its constituent `{T}`s."
)
```

```go
const (
	DefaultFiltNonNilsDocComment = "{N} returns only the non-`nil` `{T}` objects contained in `{this}`."
	DefaultFiltFuncDocComment    = "{N} returns only the `{T}` objects contained in `{this}` that satisfy the specified `{ok}` predicate."
	DefaultFiltByDocComment      = "{N} returns {what} `{T}` object(s) encountered in `{this}` whose `{member}` member succeeds for the specified value(s)."
)
```

```go
const (
	DefaultIndexMethodName           = "Index"
	DefaultIndicesMethodName         = "Indices"
	DefaultIndexLastMethodName       = "LastIndex"
	DefaultContainsMethodName        = "Contains"
	DefaultMethodNameSuffixPredicate = "Func"
)
```

```go
const (
	DefaultMutAppendDocComment = "{N} is a convenience (dot-accessor) short-hand for Go's built-in `append` function."
)
```

```go
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IndexOf    GentIndexMethods
		Filters    GentFilteringMethods
		Mutators   GentMutatorMethods
		Converters GentConvertMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

#### type GentConvertMethods

```go
type GentConvertMethods struct {
	gent.Opts

	Fields struct {
		gent.Variant
		Named []string
	}
	ToMaps struct {
		gent.Variant
		ByFields []string
	}
}
```


#### func (*GentConvertMethods) EnableOrDisableAllVariantsAndOptionals

```go
func (this *GentConvertMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals implements
`github.com/metaleap/go-gent.IGent`.

#### func (*GentConvertMethods) GenerateTopLevelDecls

```go
func (this *GentConvertMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentFilteringMethods

```go
type GentFilteringMethods struct {
	gent.Opts

	NonNils gent.Variant
	Func    gent.Variant
	By      struct {
		gent.Variation
		Fields  []string
		Methods []NamedTyped
	}
}
```


#### func (*GentFilteringMethods) EnableOrDisableAllVariantsAndOptionals

```go
func (this *GentFilteringMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals implements
`github.com/metaleap/go-gent.IGent`.

#### func (*GentFilteringMethods) GenerateTopLevelDecls

```go
func (this *GentFilteringMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentIndexMethods

```go
type GentIndexMethods struct {
	gent.Opts

	IndexOf struct {
		IndexMethodOpts
		Variadic bool
	}

	// `Disabled` in `Gents.IndexOf` by default
	IndexLast struct {
		IndexMethodOpts
		Variadic bool
	}

	// `Disabled` in `Gents.IndexOf` by default
	IndicesOf struct {
		IndexMethodOpts
		ResultsCapFactor uint
	}

	// `Disabled` in `Gents.IndexOf` by default
	Contains struct {
		IndexMethodOpts
		VariadicAny bool
		VariadicAll bool
	}
}
```


#### func (*GentIndexMethods) EnableOrDisableAllVariantsAndOptionals

```go
func (this *GentIndexMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals implements
`github.com/metaleap/go-gent.IGent`.

#### func (*GentIndexMethods) GenerateTopLevelDecls

```go
func (this *GentIndexMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentMutatorMethods

```go
type GentMutatorMethods struct {
	gent.Opts

	Append gent.Variant
}
```


#### func (*GentMutatorMethods) EnableOrDisableAllVariantsAndOptionals

```go
func (this *GentMutatorMethods) EnableOrDisableAllVariantsAndOptionals(enabled bool)
```
EnableOrDisableAllVariantsAndOptionals implements
`github.com/metaleap/go-gent.IGent`.

#### func (*GentMutatorMethods) GenerateTopLevelDecls

```go
func (this *GentMutatorMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type IndexMethodOpts

```go
type IndexMethodOpts struct {
	gent.Variation
	Predicate gent.Variant
}
```
