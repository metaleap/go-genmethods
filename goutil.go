package gent

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/go-leap/dev/go/gen"
)

func goAstTypeExprSansParens(expr ast.Expr) ast.Expr {
	for par, is := expr.(*ast.ParenExpr); is; par, is = expr.(*ast.ParenExpr) {
		expr = par.X
	}
	switch x := expr.(type) {
	case *ast.StarExpr:
		x.X = goAstTypeExprSansParens(x.X)
	case *ast.SelectorExpr:
		x.X = goAstTypeExprSansParens(x.X)
	case *ast.ArrayType:
		x.Elt = goAstTypeExprSansParens(x.Elt)
	case *ast.ChanType:
		x.Value = goAstTypeExprSansParens(x.Value)
	case *ast.MapType:
		x.Key, x.Value = goAstTypeExprSansParens(x.Key), goAstTypeExprSansParens(x.Value)
	}
	return expr
}

func goAstTypeExprToGenTypeRef(expr ast.Expr) *udevgogen.TypeRef {
	switch tx := expr.(type) {
	case *ast.Ident:
		return udevgogen.TrNamed("", tx.Name)
	case *ast.SelectorExpr:
		return udevgogen.TrNamed(tx.Sel.Name, tx.X.(*ast.Ident).Name)
	case *ast.StarExpr:
		return udevgogen.TrPtr(goAstTypeExprToGenTypeRef(tx.X))
	case *ast.ArrayType:
		switch l := tx.Len.(type) {
		case *ast.BasicLit:
			fixedlen, err := strconv.ParseUint(l.Value, 0, 64)
			if err != nil {
				panic(err)
			}
			return udevgogen.TrArray(fixedlen, goAstTypeExprToGenTypeRef(tx.Elt))
		default:
			return udevgogen.TrSlice(goAstTypeExprToGenTypeRef(tx.Elt))
		}
	case *ast.Ellipsis:
		sl := udevgogen.TrSlice(goAstTypeExprToGenTypeRef(tx.Elt))
		sl.ArrOrSliceIsEllipsis = true
		return sl
	case *ast.MapType:
		return udevgogen.TrMap(goAstTypeExprToGenTypeRef(tx.Key), goAstTypeExprToGenTypeRef(tx.Value))
	case *ast.FuncType:
		var tdfn udevgogen.TypeFunc
		if tx.Params != nil {
			for _, fld := range tx.Params.List {
				tdfn.Args.Add("", goAstTypeExprToGenTypeRef(fld.Type))
			}
		}
		if tx.Results != nil {
			for _, fld := range tx.Results.List {
				tdfn.Rets.Add("", goAstTypeExprToGenTypeRef(fld.Type))
			}
		}
		return udevgogen.TrFunc(&tdfn)
	case *ast.InterfaceType:
		if tx.Incomplete {
			panic("interface-type methods list incomplete: investigate to handle!")
		}
		var tdi udevgogen.TypeInterface
		if tx.Methods != nil {
			for _, fld := range tx.Methods.List {
				var fldname string
				if len(fld.Names) == 1 {
					fldname = fld.Names[0].Name
				}
				if len(fld.Names) == 0 {
					tdi.Embeds = append(tdi.Embeds, goAstTypeExprToGenTypeRef(fld.Type))
				} else {
					tdi.Methods.Add(fldname, goAstTypeExprToGenTypeRef(fld.Type))
				}
			}
		}
		return udevgogen.TrInterface(&tdi)
	case *ast.StructType:
		if tx.Incomplete {
			panic("struct-type fields list incomplete: investigate to handle!")
		}
		var tds udevgogen.TypeStruct
		if tx.Fields != nil {
			for _, fld := range tx.Fields.List {
				var fldname, fldtag string
				if l := len(fld.Names); l == 1 {
					fldname = fld.Names[0].Name
				} else if l > 1 {
					panic("multiple struct field names? investigate to handle!")
				}
				if fld.Tag != nil && fld.Tag.Kind == token.STRING {
					fldtag, _ = strconv.Unquote(fld.Tag.Value)
				}
				tds.Fields = append(tds.Fields, udevgogen.TdStructFld(fldname, goAstTypeExprToGenTypeRef(fld.Type), goStructFieldTagsTryParse(fldtag)))
			}
		}
		return udevgogen.TrStruct(&tds)
	case *ast.ChanType:
		return udevgogen.TrChan(tx.Dir == ast.RECV, tx.Dir == ast.SEND, goAstTypeExprToGenTypeRef(tx.Value))
	default:
		panic(tx)
	}
}

func goStructFieldTagsTryParse(tagLiteral string) (tags map[string]string) {
	tags = map[string]string{}
	for len(tagLiteral) > 0 {
		poscolonquote := strings.Index(tagLiteral, `:"`)
		if poscolonquote <= 0 {
			break
		}

		tagname := tagLiteral[:poscolonquote]
		if tagname = strings.TrimSpace(tagname); tagname == "" {
			break
		}

		tagvalpos, tagvallen := poscolonquote+2, -1
		for i, r := range tagLiteral[tagvalpos:] {
			if r == '"' && (i == 0 || tagLiteral[i-1] != '\\') {
				tagvallen = i
				break
			}
		}
		if tagvallen < 0 {
			break
		}

		tagval := tagLiteral[tagvalpos : tagvalpos+tagvallen]
		tags[tagname], tagLiteral = tagval, tagLiteral[tagvalpos+tagvallen+1:]
	}
	if tagLiteral = strings.TrimSpace(tagLiteral); tagLiteral != "" {
		tags[""] = tagLiteral
	} else if len(tags) == 0 {
		tags = nil
	}
	return
}
