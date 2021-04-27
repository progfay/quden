package gorilla

import (
	"go/ast"
)

var muxInstance instance = &mux{}

type mux struct{}

func (m *mux) Call(name string, args ...ast.Expr) instance {
	switch name {
	case "NewRouter":
		return m.NewRouter(args...)

	default:
		return m
	}
}

// HandleFunc emulate behavior of mux.NewRouter
// Ref. https://pkg.go.dev/github.com/gorilla/mux#NewRouter
func (m *mux) NewRouter(args ...ast.Expr) *Router {
	if len(args) != 0 {
		return nil
	}

	return &Router{
		matchers: []matcher{},
		subs:     []instance{},
	}
}
