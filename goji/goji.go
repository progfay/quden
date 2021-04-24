package goji

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/progfay/quden/util"
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

type visitor struct {
	endpoints []util.Endpoint
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

	x, ok := selectorExpr.X.(*ast.Ident)
	if !ok || x.Name != "pat" {
		return v
	}

	name := selectorExpr.Sel.Name
	if !isRegisterMethod(name) {
		return v
	}

	v.endpoints = append(v.endpoints, util.NewEndpoint(strings.ToUpper(name), path, path))

	return v
}

type goji struct{}

func (goji) MatchImportPath(path string) bool {
	return path == "goji.io/pat"
}

func (goji) Extract(node ast.Node) []util.Endpoint {
	var v visitor
	ast.Walk(&v, node)
	// Sort
	return v.endpoints
}

func New() util.Framework {
	return goji{}
}
