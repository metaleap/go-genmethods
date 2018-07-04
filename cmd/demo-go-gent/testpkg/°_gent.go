package genttest

// DO NOT EDIT: code generated with `demo-go-gent` using `github.com/metaleap/go-gent`

import (
	pkg__strconv "strconv"
	pkg__strings "strings"
)

// Iszero returns whether the value of this `Num` equals `zero`.
func (this Num) Iszero() (r bool) { r = this == zero; return }

// Isone returns whether the value of this `Num` equals `one`.
func (this Num) Isone() (r bool) { r = this == one; return }

// Istwo returns whether the value of this `Num` equals `two`.
func (this Num) Istwo() (r bool) { r = this == two; return }

// Isthree returns whether the value of this `Num` equals `three`.
func (this Num) Isthree() (r bool) { r = this == three; return }

// Valid returns whether the value of this `Num` is between `zero` (inclusive) and `three` (inclusive).
func (this Num) Valid() (r bool) { r = (this >= zero) && (this <= three); return }

// WellknownNums returns the `names` and `values` of all 4 well-known `Num` enumerants.
func WellknownNums() (names []string, values []Num) {
	names, values = []string{"zero", "one", "two", "three"}, []Num{zero, one, two, three}
	return
}

// String implements the `fmt.Stringer` interface.
func (this Num) String() (r string) {
	if (this < zero) || (this > three) {
		goto formatNum
	}
	switch this {
	case zero:
		r = "Zero"
	case one:
		r = "One"
	case two:
		r = "Two"
	case three:
		r = "Three"
	default:
		goto formatNum
	}
	return
formatNum:
	r = pkg__strconv.Itoa((int)(this))
	return
}

// NumFromString returns the `Num` represented by `s` (as returned by `Num.String`, but case-insensitively), or an `error` if none exists.
func NumFromString(s string) (this Num, err error) {
	if (len(s) < 3) || (len(s) > 5) {
		goto tryParseNum
	}
	switch {
	case pkg__strings.EqualFold(s, "Zero"):
		this = zero
	case pkg__strings.EqualFold(s, "One"):
		this = one
	case pkg__strings.EqualFold(s, "Two"):
		this = two
	case pkg__strings.EqualFold(s, "Three"):
		this = three
	default:
		goto tryParseNum
	}
	return
tryParseNum:
	var v int
	v, err = pkg__strconv.Atoi(s)
	if err == nil {
		this = (Num)(v)
	}
	return
}

// NumFromStringOr is like `NumFromString` but returns `fallback` for unrecognized inputs.
func NumFromStringOr(s string, fallback Num) (this Num) {
	maybeNum, err := NumFromString(s)
	if err == nil {
		this = maybeNum
	} else {
		this = fallback
	}
	return
}

// GoString implements the `fmt.GoStringer` interface.
func (this Num) GoString() (r string) {
	if (this < zero) || (this > three) {
		goto formatNum
	}
	switch this {
	case zero:
		r = "zero"
	case one:
		r = "one"
	case two:
		r = "two"
	case three:
		r = "three"
	default:
		goto formatNum
	}
	return
formatNum:
	r = pkg__strconv.Itoa((int)(this))
	return
}

// NumFromGoString returns the `Num` represented by `s` (as returned by `Num.GoString`, and case-sensitively), or an `error` if none exists.
func NumFromGoString(s string) (this Num, err error) {
	if (len(s) < 3) || (len(s) > 5) {
		goto tryParseNum
	}
	switch s {
	case "zero":
		this = zero
	case "one":
		this = one
	case "two":
		this = two
	case "three":
		this = three
	default:
		goto tryParseNum
	}
	return
tryParseNum:
	var v int
	v, err = pkg__strconv.Atoi(s)
	if err == nil {
		this = (Num)(v)
	}
	return
}

// NumFromGoStringOr is like `NumFromGoString` but returns `fallback` for unrecognized inputs.
func NumFromGoStringOr(s string, fallback Num) (this Num) {
	maybeNum, err := NumFromGoString(s)
	if err == nil {
		this = maybeNum
	} else {
		this = fallback
	}
	return
}

func (this complex384) Index(v complex128) (r int) {
	for i := range this {
		if this[i] == v {
			r = i
			return
		}
	}
	r = -1
	return
}

func (this complex384) IndexFunc(ok func(complex128) bool) (r int) {
	for i := range this {
		if ok(this[i]) {
			r = i
			return
		}
	}
	r = -1
	return
}

func (this complex384) LastIndex(v complex128) (r int) {
	for i := len(this) - 1; i > -1; i-- {
		if this[i] == v {
			r = i
			return
		}
	}
	r = -1
	return
}

func (this complex384) LastIndexFunc(ok func(complex128) bool) (r int) {
	for i := len(this) - 1; i > -1; i-- {
		if ok(this[i]) {
			r = i
			return
		}
	}
	r = -1
	return
}

func (this complex384) Indices(v complex128) (r []int) {
	for i := range this {
		if this[i] == v {
			r = append(r, i)
		}
	}
	return
}

func (this complex384) IndicesFunc(ok func(complex128) bool) (r []int) {
	for i := range this {
		if ok(this[i]) {
			r = append(r, i)
		}
	}
	return
}
