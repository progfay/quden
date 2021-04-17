package visitor

import (
	"go/ast"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/endpoint"
)

type visitor struct {
	endpoints []endpoint.Endpoint
}

func New() *visitor {
	return &visitor{
		endpoints: make([]endpoint.Endpoint, 0),
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	endpoint := echo.NodeToEndpoint(node)
	if endpoint != nil {
		v.endpoints = append(v.endpoints, *endpoint)
	}

	return v
}

func (v *visitor) GetEndpoints() []endpoint.Endpoint {
	return v.endpoints
}
