package echo

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/progfay/quden/util"
)

var registerMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

func isRegisterMethod(name string) bool {
	for _, method := range registerMethods {
		if method == name {
			return true
		}
	}
	return false
}

type visitor struct {
	utils []util.Endpoint
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
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

	method := selectorExpr.Sel.Name
	if !isRegisterMethod(method) {
		return v
	}

	v.utils = append(v.utils, util.NewEndpoint(method+path, method+path))

	return v
}

type echo struct{}

func (echo) MatchImportPath(path string) bool {
	return strings.HasPrefix(path, "github.com/labstack/echo/")
}

func (echo) Extract(node ast.Node) []util.Endpoint {
	var v visitor
	ast.Walk(&v, node)
	// Sort
	return v.utils
}

func New() util.Framework {
	return echo{}
}
