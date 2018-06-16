package gents

type Pkg struct {
	Gens map[string]Gen
}

func New() *Pkg {
	pkg := &Pkg{Gens: map[string]Gen{}}
	return pkg
}

func (this *Pkg) Load() (err error) {
	return
}

func (this *Pkg) Generate() (err error) {
	return
}
