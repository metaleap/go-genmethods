package genttest

// DO NOT EDIT: code generated with `demo-go-gent` using `github.com/metaleap/go-gent`

func (this complex384) Indices(eq complex128) (r []int) {
	for i := range this {
		if this[i] == eq {
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
