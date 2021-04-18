package visitor

import (
	"go/ast"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/endpoint"
	"github.com/progfay/quden/framework"
)

var frameworks = []*framework.Framework{echo.New()}

type visitor struct {
	endpoints      []endpoint.Endpoint
	nodeConverters []framework.NodeConverter
}

func New() *visitor {
	return &visitor{
		endpoints: make([]endpoint.Endpoint, 0),
	}
}

func (v *visitor) AddNodeConverter(nodeConverter framework.NodeConverter) {
	v.nodeConverters = append(v.nodeConverters, nodeConverter)
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	for _, nodeConverter := range v.nodeConverters {
		endpoint := nodeConverter.ToEndpoint(node)
		if endpoint != nil {
			v.endpoints = append(v.endpoints, *endpoint)
			return v
		}
	}

	return v
}

func (v *visitor) GetEndpoints() []endpoint.Endpoint {
	return v.endpoints
}
