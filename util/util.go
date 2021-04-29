package util

import (
	"fmt"
	"go/ast"
)

const bundleTemplate = `[[bundle]]
regexp = "%s"
name = "%s"`

type Endpoint struct {
	Name   string
	RegExp string
}

func NewEndpoint(name, regexp string) Endpoint {
	return Endpoint{
		Name:   name,
		RegExp: regexp,
	}
}

func (e *Endpoint) Format() string {
	return fmt.Sprintf(bundleTemplate, e.RegExp, e.Name)
}

func (e *Endpoint) String() string {
	return e.Name
}

type Framework interface {
	MatchImportPath(path string) bool
	Extract(node ast.Node) []Endpoint
}
