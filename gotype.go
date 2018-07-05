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

	// Expr is whatever underlying-type this type-decl represents, that is:
	// of the original `type foo bar` or `type foo = bar` declaration,
	// this `Type` is the `foo` identity and its `Expr` captures the `bar`.
	Expr struct {
		// original AST's type-decl's `Expr` (stripped of any&all `ParenExpr`s)
		AstExpr ast.Expr
		// a code-gen `TypeRef` to this `Type` decl's underlying-type
		GenRef *udevgogen.TypeRef
	}

	// commonly useful code-gen values prepared for this `Type`
	G struct {
		// a type-ref to this `Type`
		T *udevgogen.TypeRef
		// a type-ref to pointer-to-`Type` (think ª for addr)
		Tª *udevgogen.TypeRef
		// a type-ref to slice-of-`Type`
		Ts *udevgogen.TypeRef
		// a type-ref to slice-of-pointers-to-`Type`
		Tªs *udevgogen.TypeRef
		// Name="this" and Type=.G.T
		This udevgogen.NamedTyped
		// Name="this" and Type=.G.Tª
		Thisª udevgogen.NamedTyped
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
					t.G.T, t.Expr.AstExpr = udevgogen.TrNamed("", t.Name), goAstTypeExprSansParens(tdecl.Type)
					t.G.Tª, t.G.Ts = udevgogen.TrPtr(t.G.T), udevgogen.TrSlice(t.G.T)
					t.G.Tªs, t.G.This, t.G.Thisª = udevgogen.TrSlice(t.G.Tª), udevgogen.Vars.This.OfType(t.G.T), udevgogen.Vars.This.OfType(udevgogen.TrPtr(t.G.T))
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
		if t.Expr.GenRef = goAstTypeExpr2GenTypeRef(t.Expr.AstExpr); t.Expr.GenRef == nil {
			panic(t.Expr.AstExpr)
		}
	}
	for _, t := range this.Types {
		t.setPotentiallyEnumish()
	}
}

func (this *Type) setPotentiallyEnumish() {
	if this.Enumish.BaseType = ""; this.Expr.GenRef.Named.PkgName == "" && len(this.Enumish.ConstNames) > 0 {
		switch this.Expr.GenRef.Named.TypeName {
		case "int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
			this.Enumish.BaseType = this.Expr.GenRef.Named.TypeName
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
	return this.Expr.GenRef.ArrOrSlice.Of != nil
}
