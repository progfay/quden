package gorilla

import (
	"go/ast"
)

type instanceType = int

const (
	MUX instanceType = iota
	ROUTER
	ROUTE
)

type instance interface {
	Call(name string, args ...ast.Expr) instance
}
