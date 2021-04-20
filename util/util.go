package util

import (
	"fmt"
	"go/ast"
)

type Endpoint struct {
	Method string
	Path   string
	RegExp string
}

func NewEndpoint(method, path, regexp string) Endpoint {
	return Endpoint{
		Method: method,
		Path:   path,
		RegExp: regexp,
	}
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s %s", e.Method, e.Path)
}

type Framework interface {
	MatchImportPath(path string) bool
	Extract(node ast.Node) []Endpoint
}
