package visitor

import (
	"go/ast"
	"io"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/endpoint"
)

type visitor struct {
	w         io.Writer
	endpoints []*endpoint.Endpoint
}

func New(w io.Writer) *visitor {
	return &visitor{
		w:         w,
		endpoints: make([]*endpoint.Endpoint, 0),
	}
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	endpoint := echo.NodeToEndpoint(node)
	if endpoint != nil {
		v.endpoints = append(v.endpoints, endpoint)
	}

	return v
}
