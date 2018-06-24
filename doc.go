// Package gent offers a _Golang code-**gen t**oolkit_; philosophy being:
//
// - _"your package's type-defs. your pluggable custom code-gen logics (+ many built-in ones), tuned via your struct-field tags. one `go generate` call."_
//
// The design idea is that your codegen programs remains your own `main`
// packages written by you, but importing `gent` keeps them short and
// high-level: fast and simple to write, iterate, maintain over time.
// Furthermore (unlike unwieldy config-file-formats or 100s-of-cmd-args)
// this approach grants Turing-complete control over fine-tuning the code-gen
// flow to only generate what's truly needed, rather than "every possible func
// for every possible type-def", to minimize both code-gen and compilation times.
//
// Focus at the beginning is strictly on generating `func`s and methods for a
// package's _existing type-defs_, **not** generating type-defs such as `struct`s.
//
// For building the AST of the to-be-emitted Go source file:
//
// - `gent` relies on my `github.com/go-leap/dev/go/gen` package
//
// - and so do the built-in code-gens under `github.com/metaleap/go-gent/gents/...`,
//
// - but your custom `gent.IGent` implementers are free to prefer other
// approaches (such as `text/template` or `github.com/dave/jennifer` or
// hard-coded string-building or other) by having their `GenerateTopLevelDecls`
// implementation return a `github.com/go-leap/dev/go/gen.SynRaw`-typed byte-array.
//
// Very WIP: more comprehensive readme / package docs to come.
package gent
