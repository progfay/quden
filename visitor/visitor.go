package visitor

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"strconv"
	"strings"
)

type visitor struct{
	w io.Writer
}

func New(w io.Writer) *visitor {
	return &visitor{w: w}
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	if len(callExpr.Args) < 1 {
		return v
	}

	firstArg, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok || firstArg.Kind != token.STRING {
		return v
	}
	path, err := strconv.Unquote(firstArg.Value)
	if err != nil || !strings.HasPrefix(path, "/") {
		return v
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	name := selectorExpr.Sel.Name

	fmt.Fprintf(v.w, "%s %s\n", strings.ToUpper(name), path)
	return v
}
