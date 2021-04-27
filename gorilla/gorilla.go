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
	// Sort

	var endpoints []util.Endpoint
	for _, inst := range v.instanceMap {
		route, ok := inst.(*Route)
		if !ok || !route.isHandled {
			continue
		}
		log.Printf("%#v\n", route)
	}
	return endpoints
}

func New() util.Framework {
	return gorilla{}
}
