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

func (this Types) Struct(name string) *Type {
	if t := this.Named(name); t != nil && t.Ast.TStruct != nil {
		return t
	}
	return nil
}

type Type struct {
	pkg *Pkg

	Name  string
	Decl  *ast.TypeSpec
	Alias bool

	Ast struct {
		Named      *ast.Ident
		Imported   *ast.SelectorExpr
		Ptr        *ast.StarExpr
		TArrOrSl   *ast.ArrayType
		TChan      *ast.ChanType
		TFunc      *ast.FuncType
		TInterface *ast.InterfaceType
		TMap       *ast.MapType
		TStruct    *ast.StructType
	}

	Enumish struct {
		Potentially bool
		ConstNames  []string
	}

	CodeGen struct {
		MethodRecvVal udevgogen.NamedTyped
		MethodRecvPtr udevgogen.NamedTyped
	}
}

func (this *Pkg) load_Types(goFile *ast.File) {
	for _, topleveldecl := range goFile.Decls {
		if somedecl, _ := topleveldecl.(*ast.GenDecl); somedecl != nil {
			var curvaltident *ast.Ident
			for _, spec := range somedecl.Specs {
				if tdecl, _ := spec.(*ast.TypeSpec); tdecl != nil && tdecl.Name != nil && tdecl.Name.Name != "" && tdecl.Type != nil {
					tdx, t := goAstTypeExprSansParens(tdecl.Type), &Type{pkg: this, Name: tdecl.Name.Name, Decl: tdecl, Alias: tdecl.Assign.IsValid()}
					t.CodeGen.MethodRecvVal, t.CodeGen.MethodRecvPtr = udevgogen.V.This.Typed(udevgogen.TrNamed("", t.Name)), udevgogen.V.This.Typed(udevgogen.TrPtr(udevgogen.TrNamed("", t.Name)))
					this.Types.Add(t)

					switch tdeclt := tdx.(type) {
					case *ast.Ident:
						t.Ast.Named = tdeclt
						t.setPotentiallyEnumish()
					case *ast.StarExpr:
						t.Ast.Ptr = tdeclt
					case *ast.SelectorExpr:
						t.Ast.Imported = tdeclt
					case *ast.ArrayType:
						t.Ast.TArrOrSl = tdeclt
					case *ast.ChanType:
						t.Ast.TChan = tdeclt
					case *ast.FuncType:
						t.Ast.TFunc = tdeclt
					case *ast.InterfaceType:
						t.Ast.TInterface = tdeclt
					case *ast.MapType:
						t.Ast.TMap = tdeclt
					case *ast.StructType:
						t.Ast.TStruct = tdeclt
					default:
						panic(tdeclt)
					}
				} else if cdecl, _ := spec.(*ast.ValueSpec); somedecl.Tok == token.CONST && cdecl != nil && len(cdecl.Names) == 1 {
					if cdecl.Type != nil {
						curvaltident, _ = cdecl.Type.(*ast.Ident)
					}
					if curvaltident != nil {
						if tnamed := this.Types.Named(curvaltident.Name); tnamed != nil && tnamed.Enumish.Potentially {
							tnamed.Enumish.ConstNames = append(tnamed.Enumish.ConstNames, cdecl.Names[0].Name)
						}
					}
				}
			}
		}
	}
}

func (this *Type) setPotentiallyEnumish() {
	if this.Enumish.Potentially = false; this.Ast.Named != nil && !this.Ast.Named.IsExported() {
		switch this.Ast.Named.Name {
		case "int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "byte", "rune":
			this.Enumish.Potentially = true
		}
	}
}
