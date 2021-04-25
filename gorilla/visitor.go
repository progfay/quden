package gorilla

import (
	"go/ast"
)

type visitor struct {
	instanceMap map[*ast.CallExpr]instance
}

func (v *visitor) visit(callExpr *ast.CallExpr) instance {
	if callExpr == nil {
		return nil
	}

	if inst, ok := v.instanceMap[callExpr]; ok {
		return inst
	}

	if inst, ok := v.instanceMap[callExpr]; ok {
		return inst
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	switch x := selectorExpr.X.(type) {
	case *ast.Ident:
		if x.Name == "mux" && selectorExpr.Sel.Name == "NewRouter" {
			inst := muxInstance.Call("NewRouter", callExpr.Args...)
			v.instanceMap[callExpr] = inst
			return inst
		}

	case *ast.CallExpr:
		if inst := v.visit(x); inst != nil {
			ret := inst.Call(selectorExpr.Sel.Name, callExpr.Args...)
			if ret != nil && ret != inst {
				v.instanceMap[callExpr] = ret
			}
			return ret
		}
	}

	ident, ok := selectorExpr.X.(*ast.Ident)
	if !ok {
		return nil
	}

	if ident.Obj == nil {
		return nil
	}
	astStmt, ok := ident.Obj.Decl.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	if len(astStmt.Rhs) == 0 {
		return nil
	}
	rh, ok := astStmt.Rhs[0].(*ast.CallExpr)
	if !ok {
		return nil
	}
	if inst := v.visit(rh); inst != nil {
		ret := inst.Call(selectorExpr.Sel.Name, callExpr.Args...)
		if ret != nil && ret != inst {
			v.instanceMap[callExpr] = ret
		}
		return ret
	}
	return nil
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.AssignStmt:
		for _, rh := range node.Rhs {
			if callExpr, ok := rh.(*ast.CallExpr); ok {
				v.visit(callExpr)
			}
		}

	case *ast.ExprStmt:
		if callExpr, ok := node.X.(*ast.CallExpr); ok {
			v.visit(callExpr)
		}

	default:
		return v
	}

	return nil
}
