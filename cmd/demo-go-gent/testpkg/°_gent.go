package genttest

// DO NOT EDIT: code generated with `demo-go-gent` using `github.com/metaleap/go-gent`

func (this complex384) Index(eq complex128) (r int) {
	for i := range this {
		if this[i] == eq {
			r = i
			return
		}
	}
	r = -1
	return
}

func (this complex384) Contains(eq complex128) (r bool) { return }
