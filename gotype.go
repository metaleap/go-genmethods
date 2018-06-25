package gent

import (
	"go/ast"
	"go/token"

	"github.com/go-leap/dev/go/gen"
)

type Types []*Type

func (this *Types) Add(t *Type) {
	*this = append(*this, t)
}

func (this Types) Named(name string) *Type {
	for _, t := range this {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type Type struct {
	Pkg *Pkg

	Name  string
	Alias bool

	Underlying struct {
		AstExpr ast.Expr
		GenRef  *udevgogen.TypeRef
	}

	Gen struct {
		TVal    *udevgogen.TypeRef
		TPtr    *udevgogen.TypeRef
		TSl     *udevgogen.TypeRef
		ThisVal udevgogen.NamedTyped
		ThisPtr udevgogen.NamedTyped
	}

	Enumish struct {
		// expected to be builtin prim-type such as uint8, int64, int --- cases of additional indirections to be handled when they occur in practice
		BaseType string

		ConstNames []string
	}
}

func (this *Pkg) load_Types(goFile *ast.File) {
	for _, topleveldecl := range goFile.Decls {
		if somedecl, _ := topleveldecl.(*ast.GenDecl); somedecl != nil {
			var curvaltident *ast.Ident
			for _, spec := range somedecl.Specs {
				if tdecl, _ := spec.(*ast.TypeSpec); tdecl != nil && tdecl.Name != nil && tdecl.Name.Name != "" && tdecl.Type != nil {
					t := &Type{Pkg: this, Name: tdecl.Name.Name, Alias: tdecl.Assign.IsValid()}
					t.Gen.TVal, t.Underlying.AstExpr = udevgogen.TrNamed("", t.Name), goAstTypeExprSansParens(tdecl.Type)
					t.Gen.TPtr, t.Gen.TSl = udevgogen.TrPtr(t.Gen.TVal), udevgogen.TrSlice(t.Gen.TVal)
					t.Gen.ThisVal, t.Gen.ThisPtr = udevgogen.V.This.T(t.Gen.TVal), udevgogen.V.This.T(udevgogen.TrPtr(t.Gen.TVal))
					this.Types.Add(t)
				} else if cdecl, _ := spec.(*ast.ValueSpec); somedecl.Tok == token.CONST && cdecl != nil && len(cdecl.Names) == 1 {
					if cdecl.Type != nil {
						curvaltident, _ = cdecl.Type.(*ast.Ident)
					}
					if curvaltident != nil {
						if tnamed := this.Types.Named(curvaltident.Name); tnamed != nil {
							tnamed.Enumish.ConstNames = append(tnamed.Enumish.ConstNames, cdecl.Names[0].Name)
						}
					}
				}
			}
		}
	}
}

func (this *Pkg) load_PopulateTypes() {
	for _, t := range this.Types {
		if t.Underlying.GenRef = goAstTypeExprToGenTypeRef(t.Underlying.AstExpr); t.Underlying.GenRef == nil {
			panic(t.Underlying.AstExpr)
		}
	}
	for _, t := range this.Types {
		t.setPotentiallyEnumish()
	}
}

func (this *Type) setPotentiallyEnumish() {
	if this.Enumish.BaseType = ""; this.Underlying.GenRef.Named.PkgName == "" && len(this.Enumish.ConstNames) > 0 {
		switch this.Underlying.GenRef.Named.TypeName {
		case "int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
			this.Enumish.BaseType = this.Underlying.GenRef.Named.TypeName
		}
	}
	if this.Enumish.BaseType == "" {
		this.Enumish.ConstNames = nil
	}
}

func (this *Type) IsEnumish() bool {
	return this.Enumish.BaseType != "" && len(this.Enumish.ConstNames) > 0 && (this.Enumish.ConstNames[0] != "_" || len(this.Enumish.ConstNames) > 1)
}

func (this *Type) IsSliceOrArray() bool {
	return this.Underlying.GenRef.ArrOrSliceOf.Val != nil
}
