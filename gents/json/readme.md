# gentjson
--
    import "github.com/metaleap/go-gent/gents/json"


## Usage

```go
const (
	DefaultMethodNameMarshal   = "MarshalJSON"
	DefaultMethodNameUnmarshal = "UnmarshalJSON"
	DefaultDocCommentMarshal   = "{N} implements the Go standard library's `encoding/json.Marshaler` interface."
	DefaultDocCommentUnmarshal = "{N} implements the Go standard library's `encoding/json.Unmarshaler` interface."
)
```

```go
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		EnumishTypes GentEnumJsonMethods

		OtherTypes GentTypeJsonMethods

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

#### type GentEnumJsonMethods

```go
type GentEnumJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
	}
	Unmarshal struct {
		JsonMethodOpts
	}
	StringerToUse *gentenums.StringMethodOpts
}
```


#### func (*GentEnumJsonMethods) GenerateTopLevelDecls

```go
func (me *GentEnumJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentTypeJsonMethods

```go
type GentTypeJsonMethods struct {
	gent.Opts

	Marshal struct {
		JsonMethodOpts
		InitialBytesCap int
	}
	Unmarshal struct {
		JsonMethodOpts
	}
}
```


#### func (*GentTypeJsonMethods) GenerateTopLevelDecls

```go
func (me *GentTypeJsonMethods) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type JsonMethodOpts

```go
type JsonMethodOpts struct {
	gent.Variation
	MayGenFor func(*gent.Type) bool
}
```
