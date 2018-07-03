package genttest

// DO NOT EDIT: code generated with `demo-go-gent` using `github.com/metaleap/go-gent`

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
