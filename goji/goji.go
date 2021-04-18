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

//  0  *ast.CallExpr {
//  1  .  Fun: *ast.SelectorExpr {
//  2  .  .  X: *ast.Ident {
//  3  .  .  .  NamePos: -
//  4  .  .  .  Name: "pat"
//  5  .  .  }
//  6  .  .  Sel: *ast.Ident {
//  7  .  .  .  NamePos: -
//  8  .  .  .  Name: "Get"
//  9  .  .  }
// 10  .  }
// 11  .  Lparen: -
// 12  .  Args: []ast.Expr (len = 1) {
// 13  .  .  0: *ast.BasicLit {
// 14  .  .  .  ValuePos: -
// 15  .  .  .  Kind: STRING
// 16  .  .  .  Value: "\"/users/:user_id\""
// 17  .  .  }
// 18  .  }
// 19  .  Ellipsis: -
// 20  .  Rparen: -
// 21  }
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

	return endpoint.New(strings.ToUpper(name), path)
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
