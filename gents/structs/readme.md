# gentstructs
--
    import "github.com/metaleap/go-gent/gents/structs"


## Usage

```go
const (
	DefaultDocCommentGet = ""
	DefaultMethodNameGet = "StructFieldsGet"
	DefaultDocCommentSet = ""
	DefaultMethodNameSet = "StructFieldsSet"
)
```

```go
const (
	DefaultDocCommentTrav = "{N} calls `on` {nf}x: once for each field in this `{T}` with its name, its pointer, `true` if name (or embed name) begins in upper-case (else `false`), and `true` if field is an embed (else `false`)."
	DefaultMethodNameTrav = "StructFieldsTraverse"
)
```

```go
var (

	// These "default `IGent`s" are a convenience offering in two ways:
	// they illustrate usage of this package's individual `IGent` implementers' fields,
	// and they allow importers their own "sane defaults" base for less-noisy tweaking.
	// They are only _initialized_ by this package, but not otherwise _used_ by it.
	Gents struct {
		StructFieldsTrav   GentStructFieldsTrav
		StructFieldsGetSet GentStructFieldsGetSet

		// contains pointers to all the above fields, in order
		All gent.Gents
	}
)
```

#### type GentStructFieldsGetSet

```go
type GentStructFieldsGetSet struct {
	gent.Opts

	Getter struct {
		gent.Variation
		ReturnsPtrInsteadOfVal bool
	}
	Setter gent.Variation
}
```


#### func (*GentStructFieldsGetSet) GenerateTopLevelDecls

```go
func (me *GentStructFieldsGetSet) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.

#### type GentStructFieldsTrav

```go
type GentStructFieldsTrav struct {
	gent.Opts

	DocComment      gent.Str
	MethodName      string
	MayIncludeField func(*SynStructField) bool
}
```


#### func (*GentStructFieldsTrav) GenerateTopLevelDecls

```go
func (me *GentStructFieldsTrav) GenerateTopLevelDecls(ctx *gent.Ctx, t *gent.Type) (yield Syns)
```
GenerateTopLevelDecls implements `github.com/metaleap/go-gent.IGent`.
