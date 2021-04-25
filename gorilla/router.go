package gorilla

import (
	"go/ast"
)

type Router struct {
	matchers []matcher

	subs   []*Router
	routes []*Route
}

func (router *Router) Call(name string, args ...ast.Expr) instance {
	switch name {
	case "HandleFunc":
		return router.HandleFunc(args...)

	case "PathPrefix":
		return router.HandleFunc(args...)

	case "Path":
		return router.Path(args...)

	case "Methods":
		return router.Methods(args...)

	case "NewRoute":
		return router.NewRoute(args...)

	default:
		return router
	}
}

// HandleFunc emulate behavior of mux.Router.HandleFunc
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.HandleFunc
func (router *Router) HandleFunc(args ...ast.Expr) *Route {
	if router == nil || len(args) != 2 {
		return nil
	}

	return router.NewRoute().Path(args[0]).HandlerFunc(args[1])
}

// PathPrefix emulate behavior of mux.Router.PathPrefix
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.PathPrefix
func (router *Router) PathPrefix(args ...ast.Expr) *Route {
	if router == nil || len(args) != 1 {
		return nil
	}

	return router.NewRoute().PathPrefix(args[0])
}

// Path emulate behavior of mux.Router.PathPrefix
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.Path
func (router *Router) Path(args ...ast.Expr) *Route {
	if router == nil || len(args) != 1 {
		return nil
	}

	return router.NewRoute().Path(args[0])
}

// Methods emulate behavior of mux.Router.Methods
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.Methods
func (router *Router) Methods(args ...ast.Expr) *Route {
	if router == nil {
		return nil
	}

	return router.NewRoute().Methods(args...)
}

// Subrouter emulate behavior of mux.Router.NewRoute
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.NewRoute
func (router *Router) NewRoute(args ...ast.Expr) *Route {
	if len(args) != 0 {
		return nil
	}

	matchers := make([]matcher, len(router.matchers))
	copy(matchers, router.matchers)

	route := &Route{
		router:    router,
		matchers:  matchers,
		isHandled: false,
	}

	router.routes = append(router.routes, route)
	return route
}
