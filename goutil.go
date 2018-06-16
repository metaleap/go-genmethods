package gent

import (
	"go/ast"
)

func goAstExprSansParens(expr ast.Expr) ast.Expr {
	for par, is := expr.(*ast.ParenExpr); is; par, is = expr.(*ast.ParenExpr) {
		expr = par.X
	}
	switch x := expr.(type) {
	case *ast.StarExpr:
		x.X = goAstExprSansParens(x.X)
	case *ast.SelectorExpr:
		x.X = goAstExprSansParens(x.X)
	case *ast.ArrayType:
		x.Elt = goAstExprSansParens(x.Elt)
	case *ast.ChanType:
		x.Value = goAstExprSansParens(x.Value)
	case *ast.MapType:
		x.Key, x.Value = goAstExprSansParens(x.Key), goAstExprSansParens(x.Value)
	}
	return expr
}
