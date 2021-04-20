package framework

import (
	"go/ast"

	"github.com/progfay/quden/endpoint"
)

type Framework interface {
	MatchImportPath(path string) bool
	Extract(node ast.Node) []endpoint.Endpoint
}
