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
	v := visitor{
		instanceMap: make(map[*ast.CallExpr]instance),
	}
	ast.Walk(&v, node)
	var endpoints []util.Endpoint
	dfs(v.entrypoint.subs)
	return endpoints
}

func dfs(insts []instance) {
	for _, inst := range insts {
		switch inst := inst.(type) {
		case *mux:
			continue

		case *Route:
			if inst.isHandled {
				log.Printf("%#v\n", inst)
			}

		case *Router:
			dfs(inst.subs)
		}
	}
}

func New() util.Framework {
	return gorilla{}
}
