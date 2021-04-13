package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

func formatBundle(regexp, name string) string {
	return fmt.Sprintf("[[bundle]]\nregexp = %q\nname = %q", regexp, name)
}

type Visitor struct{}

func (v Visitor) Visit(node ast.Node) ast.Visitor {
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

	fmt.Printf("%s %s\n", strings.ToUpper(name), path)
	return v
}

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "./example/main.go", nil, parser.Mode(0))
	v := Visitor{}

	ast.Walk(v, f)
}
