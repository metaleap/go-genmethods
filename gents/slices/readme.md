# gentslices
--
    import "github.com/metaleap/go-gent/gents/slices"


## Usage

```go
var (
	// These "default `IGent`s" are convenience offerings in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		IndexOf GentIndexMethods

		// contains pointers to all the above fields, in order
		All []gent.IGent
	}
)
```

#### type GentIndexMethods

```go
type GentIndexMethods struct {
	gent.Opts

	IndexOf struct {
		IndexMethodOpts
		gent.Variadic
	}
	IndexLast struct {
		IndexMethodOpts
		gent.Variadic
	}
	IndicesOf struct {
		IndexMethodOpts
		ResultsCapFactor uint
	}
	Contains struct {
		gent.Variant
		gent.Variadic
	}
}
```


#### func (*GentIndexMethods) GenerateTopLevelDecls

```go
func (this *GentIndexMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (decls Syns)
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
