package gent

import (
	"go/ast"
	"go/token"
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
		Potentially   bool
		ConstantNames []string
	}
}

func (this *Pkg) load_Types(goFile *ast.File) {
	for _, topleveldecl := range goFile.Decls {
		if somedecl, _ := topleveldecl.(*ast.GenDecl); somedecl != nil {
			var curvaltident *ast.Ident
			for _, spec := range somedecl.Specs {
				if tdecl, _ := spec.(*ast.TypeSpec); tdecl != nil && tdecl.Name != nil && tdecl.Name.Name != "" && tdecl.Type != nil {
					tdx, pt := goAstTypeExprSansParens(tdecl.Type), &Type{pkg: this, Name: tdecl.Name.Name, Decl: tdecl, Alias: tdecl.Assign.IsValid()}
					this.Types.Add(pt)

					switch tdeclt := tdx.(type) {
					case *ast.Ident:
						pt.Ast.Named = tdeclt
						pt.setPotentiallyEnumish()
					case *ast.StarExpr:
						pt.Ast.Ptr = tdeclt
					case *ast.SelectorExpr:
						pt.Ast.Imported = tdeclt
						println("Sel", pt.Name, tdeclt.X)
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
				} else if cdecl, _ := spec.(*ast.ValueSpec); somedecl.Tok == token.CONST && cdecl != nil && len(cdecl.Names) == 1 {
					if cdecl.Type != nil {
						curvaltident, _ = cdecl.Type.(*ast.Ident)
					}
					if curvaltident != nil {
						if tnamed := this.Types.Named(curvaltident.Name); tnamed != nil && tnamed.Enumish.Potentially {
							tnamed.Enumish.ConstantNames = append(tnamed.Enumish.ConstantNames, cdecl.Names[0].Name)
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
