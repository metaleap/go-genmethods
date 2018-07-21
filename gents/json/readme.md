# gentjson
--
    import "github.com/metaleap/go-gent/gents/json"


## Usage

```go
const (
	DefaultDocCommentMarshal   = "MarshalJSON implements the Go standard library's `encoding/json.Marshaler` interface."
	DefaultDocCommentUnmarshal = "UnmarshalJSON implements the Go standard library's `encoding/json.Unmarshaler` interface."
)
```

```go
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		Enums GentEnumJsonMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

#### type GentEnumJsonMethods

```go
type GentEnumJsonMethods struct {
	gent.Opts

	DocCommentMarshal   string
	DocCommentUnmarshal string
	StringerToUse       *gentenums.StringMethodOpts
}
```


#### func (*GentEnumJsonMethods) GenerateTopLevelDecls

```go
func (this *GentEnumJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
