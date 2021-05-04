package gorilla

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/progfay/quden/util"
)

type Route struct {
	matchers  []matcher
	isHandled bool
	router    *Router
}

func (route *Route) Call(name string, args ...ast.Expr) instance {
	if route == nil {
		return nil
	}

	switch name {
	case "PathPrefix":
		return route.PathPrefix(args...)

	case "Path":
		return route.Path(args...)

	case "Methods":
		return route.Methods(args...)

	case "HandlerFunc":
		return route.HandlerFunc(args...)

	case "Subrouter":
		return route.Subrouter(args...)

	default:
		return route
	}
}

// PathPrefix emulate behavior of mux.Route.PathPrefix
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.PathPrefix
func (route *Route) PathPrefix(args ...ast.Expr) *Route {
	if route == nil || len(args) != 1 {
		return nil
	}

	lit, ok := args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}

	pattern, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil
	}

	route.matchers = append(route.matchers, newPathPrefixMatcher(pattern))
	return route
}

// PathPrefix emulate behavior of mux.Route.Path
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.Path
func (route *Route) Path(args ...ast.Expr) *Route {
	if route == nil || len(args) != 1 {
		return nil
	}

	lit, ok := args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}

	pattern, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil
	}

	route.matchers = append(route.matchers, newPathMatcher(pattern))
	return route
}

// Methods emulate behavior of mux.Route.Methods
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.Methods
func (route *Route) Methods(args ...ast.Expr) *Route {
	if route == nil {
		return nil
	}

	methods := make([]string, len(args))

	for i, arg := range args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return nil
		}

		method, err := strconv.Unquote(lit.Value)
		if err != nil {
			return nil
		}

		methods[i] = strings.ToUpper(method)
	}

	route.matchers = append(route.matchers, newMethodsMatcher(methods))
	return route
}

// HandlerFunc emulate behavior of mux.Route.HandlerFunc
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.HandlerFunc
func (route *Route) HandlerFunc(args ...ast.Expr) *Route {
	if route == nil || len(args) != 1 {
		return nil
	}

	route.isHandled = true
	return route
}

// Subrouter emulate behavior of mux.Route.Subrouter
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.Subrouter
func (route *Route) Subrouter(args ...ast.Expr) *Router {
	if route == nil || len(args) != 0 {
		return nil
	}

	matcher := make([]matcher, len(route.matchers))
	copy(matcher, route.matchers)

	router := &Router{
		matchers: matcher,
		subs:     []instance{},
	}

	route.router.subs = append(route.router.subs, router)
	return router
}

func (route *Route) ToEndpoints() []util.Endpoint {
	if route == nil {
		return nil
	}

	art := newArtifact()
	for _, matcher := range route.matchers {
		matcher.Process(art)
	}

	name, pattern, err := parsePath(art.path)
	if err != nil {
		return nil
	}

	if !art.pathTerminated {
		pattern += `\B*`
	}
	pattern += `\b`

	if art.methodSet == nil {
		return []util.Endpoint{
			util.NewEndpoint(
				name,
				fmt.Sprintf("^[^ ]+ %s", pattern),
			),
		}
	}

	methods := make([]string, 0, len(art.methodSet))
	for method := range art.methodSet {
		methods = append(methods, method)
	}
	sort.Strings(methods)

	endpoints := make([]util.Endpoint, len(art.methodSet))
	for i, method := range methods {
		endpoints[i] = util.NewEndpoint(
			fmt.Sprintf("%s %s", method, name),
			fmt.Sprintf("^%s %s", method, pattern),
		)
	}

	return endpoints
}
