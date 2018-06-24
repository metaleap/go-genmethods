// Package gentenums provides `gent.IGent` code-gens of `func`s related to "enum-ish
// type-defs". Most of them expect and assume enum type-defs whose enumerants are
// ordered in the source such that the numerically smallest value appears first,
// the largest one last, with all enumerant `const`s appearing next to each other.
package gentenums
