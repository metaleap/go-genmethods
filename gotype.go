package gent

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/go-leap/dev/go"
	"github.com/go-leap/dev/go/gen"
)

type Types []*Type

func (me *Types) Add(t *Type) {
	*me = append(*me, t)
}

func (me Types) Named(name string) *Type {
	if name != "" {
		for _, t := range me {
			if t.Name == name {
				return t
			}
		}
	}
	return nil
}

type Type struct {
	Pkg            *Pkg
	SrcFileImports []PkgSpec

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
		// type-ref to this `Type`
		T *udevgogen.TypeRef
		// type-ref to pointer-to-`Type` (think 'ª for addr')
		Tª *udevgogen.TypeRef
		// type-ref to slice-of-`Type`
		Ts *udevgogen.TypeRef
		// type-ref to slice-of-pointers-to-`Type`
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

var miscPkgNames = map[string]string{}

func (me *Pkg) load_Types(goFile *ast.File) {
	imps := make([]PkgSpec, 0, len(goFile.Imports))
	for _, srcimp := range goFile.Imports {
		var imp PkgSpec
		if srcimp.Name != nil {
			imp.Name = srcimp.Name.Name
		}
		if srcimp.Path != nil {
			imp.ImportPath, _ = strconv.Unquote(srcimp.Path.Value)
		}
		if imp.Name != "" || imp.ImportPath != "" {
			if imp.ImportPath == "" {
				imp.ImportPath = imp.Name
			} else if imp.Name == "" {
				if imp.Name = miscPkgNames[imp.ImportPath]; imp.Name == "" {
					imp.Name = udevgo.LoadOnlyPkgNameFrom(imp.ImportPath)
					miscPkgNames[imp.ImportPath] = imp.Name
				}
			}
			imps = append(imps, imp)
		}
	}

	for _, topleveldecl := range goFile.Decls {
		if somedecl, _ := topleveldecl.(*ast.GenDecl); somedecl != nil {
			var curvaltident *ast.Ident
			for _, spec := range somedecl.Specs {
				if tdecl, _ := spec.(*ast.TypeSpec); tdecl != nil && tdecl.Name != nil && tdecl.Name.Name != "" && tdecl.Type != nil {
					t := &Type{Pkg: me, Name: tdecl.Name.Name, Alias: tdecl.Assign.IsValid(), SrcFileImports: imps}
					t.G.T, t.Expr.AstExpr = udevgogen.TFrom("", t.Name), goAstTypeExprSansParens(tdecl.Type)
					t.G.Tª, t.G.Ts = udevgogen.TPointer(t.G.T), udevgogen.TSlice(t.G.T)
					t.G.Tªs, t.G.This, t.G.Thisª = udevgogen.TSlice(t.G.Tª), udevgogen.Self.OfType(t.G.T), udevgogen.Self.OfType(udevgogen.TPointer(t.G.T))
					me.Types.Add(t)
				} else if cdecl, _ := spec.(*ast.ValueSpec); somedecl.Tok == token.CONST && cdecl != nil && len(cdecl.Names) == 1 {
					if cdecl.Type != nil {
						curvaltident, _ = cdecl.Type.(*ast.Ident)
					} else if len(cdecl.Values) > 0 {
						if call, okc := cdecl.Values[0].(*ast.CallExpr); okc {
							curvaltident, _ = call.Fun.(*ast.Ident)
						}
					}
					if curvaltident != nil {
						if tnamed := me.Types.Named(curvaltident.Name); tnamed != nil {
							tnamed.Enumish.ConstNames = append(tnamed.Enumish.ConstNames, cdecl.Names[0].Name)
						}
					}
				}
			}
		}
	}
}

func (me *Pkg) load_PopulateTypes() {
	for _, t := range me.Types {
		if t.Expr.GenRef = goAstTypeExpr2GenTypeRef(t.Expr.AstExpr); t.Expr.GenRef == nil {
			panic(t.Expr.AstExpr)
		}
	}
	for _, t := range me.Types {
		t.setPotentiallyEnumish()
	}
}

func (me *Type) setPotentiallyEnumish() {
	if me.Enumish.BaseType = ""; me.Expr.GenRef.Named.PkgName == "" && len(me.Enumish.ConstNames) > 0 {
		switch me.Expr.GenRef.Named.TypeName {
		case "int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
			me.Enumish.BaseType = me.Expr.GenRef.Named.TypeName
		}
	}
	if me.Enumish.BaseType == "" {
		me.Enumish.ConstNames = nil
	}
}

func (me *Type) IsEnumish() bool {
	return me.Enumish.BaseType != "" && len(me.Enumish.ConstNames) > 0 && (me.Enumish.ConstNames[0] != "_" || len(me.Enumish.ConstNames) > 1)
}

func (me *Type) IsArray() bool {
	return me.IsSliceOrArray() && me.Expr.GenRef.ArrOrSlice.IsFixedLen != nil
}

func (me *Type) IsSlice() bool {
	return me.IsSliceOrArray() && me.Expr.GenRef.ArrOrSlice.IsFixedLen == nil
}

func (me *Type) IsSliceOrArray() bool {
	return me.Expr.GenRef.ArrOrSlice.Of != nil
}

func (me *Type) SrcFileImportPathByName(impName string) *PkgSpec {
	if impName != "" {
		for i := range me.SrcFileImports {
			if me.SrcFileImports[i].Name == impName {
				return &me.SrcFileImports[i]
			}
		}
	}
	return nil
}
