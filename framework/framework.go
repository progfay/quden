package framework

import (
	"go/ast"

	"github.com/progfay/quden/endpoint"
)

type Framework interface {
	MatchImportPath(path string) bool
	NewNodeConverter() NodeConverter
}

type NodeConverter interface {
	ToEndpoint(node ast.Node) *endpoint.Endpoint
}
