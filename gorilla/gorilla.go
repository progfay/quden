package gorilla

import (
	"go/ast"
	"log"
	"strings"

	"github.com/progfay/quden/util"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type gorilla struct{}

func (gorilla) MatchImportPath(path string) bool {
	return strings.HasPrefix(path, "github.com/gorilla/mux")
}

func (gorilla) Extract(node ast.Node) []util.Endpoint {
	v := visitor{instanceMap: make(map[*ast.CallExpr]instance)}
	ast.Walk(&v, node)
	return dfs(v.entrypoint)
}

func dfs(inst instance) []util.Endpoint {
	switch inst := inst.(type) {
	case *mux:
		return nil

	case *Route:
		if !inst.isHandled {
			return nil
		}
		return inst.ToEndpoints()

	case *Router:
		var endpoints []util.Endpoint
		for _, sub := range inst.subs {
			endpoints = append(endpoints, dfs(sub)...)
		}
		return endpoints

	default:
		return nil
	}
}

func New() util.Framework {
	return gorilla{}
}
