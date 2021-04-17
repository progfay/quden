package echo

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/progfay/quden/endpoint"
)

var registerMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

func isRegisterMethod(name string) bool {
	for _, method := range registerMethods {
		if method == name {
			return true
		}
	}
	return false
}

func NodeToEndpoint(node ast.Node) *endpoint.Endpoint {
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

	name := strings.ToUpper(selectorExpr.Sel.Name)
	if !isRegisterMethod(name) {
		return nil
	}

	return endpoint.New(name, path)
}
