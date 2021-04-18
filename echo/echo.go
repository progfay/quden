package echo

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/progfay/quden/endpoint"
	"github.com/progfay/quden/framework"
)

var registerMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}

type echo struct {}

type converter struct {}

func New() framework.Framework {
	return echo{}
}

func isRegisterMethod(name string) bool {
	for _, method := range registerMethods {
		if method == name {
			return true
		}
	}
	return false
}

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

	name := strings.ToUpper(selectorExpr.Sel.Name)
	if !isRegisterMethod(name) {
		return nil
	}

	return endpoint.New(name, path)
}

func (echo) MatchImportPath(path string) bool {
	return strings.HasPrefix(path, "github.com/labstack/echo/")
}

func (echo) NewNodeConverter() framework.NodeConverter {
	return converter{}
}
