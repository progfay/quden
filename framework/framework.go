package framework

import (
	"go/ast"

	"github.com/progfay/quden/endpoint"
)

type Framework struct {
	MatchImportPath  func(path string) bool
	NewNodeConverter func() NodeConverter
}

type NodeConverter interface {
	ToEndpoint(node ast.Node) *endpoint.Endpoint
}
