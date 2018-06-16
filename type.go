package gent

import (
	"go/ast"
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
}

func (this *Pkg) load_Types(goFile *ast.File) {
	for _, topleveldecl := range goFile.Decls {
		if somedecl, _ := topleveldecl.(*ast.GenDecl); somedecl != nil {
			for _, spec := range somedecl.Specs {
				if tdecl, _ := spec.(*ast.TypeSpec); tdecl != nil && tdecl.Name != nil && tdecl.Name.Name != "" && tdecl.Type != nil {
					tdx, pt := goAstExprSansParens(tdecl.Type), &Type{Name: tdecl.Name.Name, Decl: tdecl, Alias: tdecl.Assign.IsValid()}
					this.Types.Add(pt)

					switch tdeclt := tdx.(type) {
					case *ast.Ident:
						pt.Ast.Named = tdeclt
					case *ast.StarExpr:
						pt.Ast.Ptr = tdeclt
					case *ast.SelectorExpr:
						pt.Ast.Imported = tdeclt
					case *ast.ArrayType:
						pt.Ast.TArrOrSl = tdeclt
					case *ast.ChanType:
						pt.Ast.TChan = tdeclt
					case *ast.FuncType:
						pt.Ast.TFunc = tdeclt
					case *ast.InterfaceType:
						pt.Ast.TInterface = tdeclt
					case *ast.MapType:
						pt.Ast.TMap = tdeclt
					case *ast.StructType:
						pt.Ast.TStruct = tdeclt
					default:
						panic(tdeclt)
					}
				}
			}
		}
	}
}
