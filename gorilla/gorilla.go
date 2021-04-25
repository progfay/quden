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
	// for _, inst := range v.instanceMap {
	// 	if !inst.isHandled {
	// 		continue
	// 	}

	// 	p := inst.GetPath()
	// 	if inst.methods == nil {
	// 		endpoints = append(endpoints, util.NewEndpoint("*", p, p))
	// 		continue
	// 	}

	// 	for _, method := range inst.methods {
	// 		endpoints = append(endpoints, util.NewEndpoint(method, p, p))
	// 	}
	// }
	return endpoints
}

func New() util.Framework {
	return gorilla{}
}
