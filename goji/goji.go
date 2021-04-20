package goji

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/progfay/quden/endpoint"
	"github.com/progfay/quden/framework"
)

var registerMethods = []string{"Get", "Head", "Post", "Put", "Delete", "CONNECT", "OPTIONS", "Patch"}

func isRegisterMethod(name string) bool {
	for _, method := range registerMethods {
		if method == name {
			return true
		}
	}
	return false
}

type converter struct{}

func (converter) ToEndpoint(node ast.Node) *endpoint.Endpoint {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}

	if len(callExpr.Args) < 1 {
		return nil
	}

	firstArg, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok || firstArg.Kind != token.STRING {
		return nil
	}

	path, err := strconv.Unquote(firstArg.Value)
	if err != nil || !strings.HasPrefix(path, "/") {
		return nil
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	x, ok := selectorExpr.X.(*ast.Ident)
	if !ok || x.Name != "pat" {
		return nil
	}

	name := selectorExpr.Sel.Name
	if !isRegisterMethod(name) {
		return nil
	}

	return endpoint.New(strings.ToUpper(name), path, path)
}

func New() *framework.Framework {
	return &framework.Framework{
		MatchImportPath: func(path string) bool {
			return path == "goji.io/pat"
		},
		NewNodeConverter: func() framework.NodeConverter {
			return converter{}
		},
	}
}
