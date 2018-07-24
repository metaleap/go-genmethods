# genttrav
--
    import "github.com/metaleap/go-gent/gents/trav"


## Usage

```go
const (
	DefaultMethodName = "TraverseFields"
)
```

```go
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		StructFieldsTrav GentStructFieldsTrav

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

#### type GentStructFieldsTrav

```go
type GentStructFieldsTrav struct {
	gent.Opts

	DocComment gent.Str
	MethodName string
}
```


#### func (*GentStructFieldsTrav) GenerateTopLevelDecls

```go
func (this *GentStructFieldsTrav) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
