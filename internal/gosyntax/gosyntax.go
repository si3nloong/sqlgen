package gosyntax

import (
	"go/ast"
	"go/types"
)

// PtrOf is to returns the pointer of the value.
func PtrOf[T any](v T) *T {
	return &v
}

// ElemOf is to get the underneath element of the pointer
// if the target is a `*ast.StarExpr`.
func ElemOf(expr ast.Expr) ast.Expr {
	for expr != nil {
		switch v := expr.(type) {
		case *ast.StarExpr:
			expr = v.X
		default:
			return v
		}
	}
	return expr
}

func UnderlyingType(t types.Type) string {
	typeStr := ""
	for t != nil {
		switch v := t.(type) {
		case *types.Pointer:
			typeStr += "*"
			t = v.Elem()
		case *types.Slice:
			typeStr += "[]"
			t = v.Elem()
		case *types.Map:
			typeStr += "map[" + UnderlyingType(v.Key()) + "]"
			t = v.Elem()
		case *types.Named:
			switch vt := t.Underlying().(type) {
			case *types.Basic:
				return typeStr + vt.String()
			default:
				return typeStr + t.String()
			}
		default:
			return typeStr + t.Underlying().String()
		}
	}
	return typeStr
}
