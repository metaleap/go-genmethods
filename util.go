package gent

import (
	"go/ast"
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
