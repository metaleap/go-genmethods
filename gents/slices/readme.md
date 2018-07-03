# gentslices
--
    import "github.com/metaleap/go-gent/gents/slices"


## Usage

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
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IndexOf GentIndexMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

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

#### type IndexMethodOpts

```go
type IndexMethodOpts struct {
	Disabled   bool
	DocComment gent.Str
	Name       string
	Predicate  gent.Variant
}
```
