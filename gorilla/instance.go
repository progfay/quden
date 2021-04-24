package gorilla

import (
	"bytes"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type instanceType = int

const (
	MUX instanceType = iota
	ROUTER
	ROUTE
)

var muxInstance = &instance{T: MUX}

type instance struct {
	T          instanceType
	superRoute *instance
	pathPrefix string
	methods    []string
	path       string
	isHandled  bool
}

func (inst *instance) Call(name string, args ...ast.Expr) *instance {
	if inst == nil {
		return nil
	}

	switch name {
	case "NewRouter":
		return inst.NewRouter(args...)

	case "HandleFunc":
		return inst.HandleFunc(args...)

	case "PathPrefix":
		return inst.PathPrefix(args...)

	case "Path":
		return inst.Path(args...)

	case "Methods":
		return inst.Methods(args...)

	case "HandlerFunc":
		return inst.HandlerFunc(args...)

	case "Subrouter":
		return inst.Subrouter(args...)

	default:
		return inst
	}
}

// HandleFunc emulate behavior of mux.NewRouter
// Ref. https://pkg.go.dev/github.com/gorilla/mux#NewRouter
func (mux *instance) NewRouter(args ...ast.Expr) *instance {
	if mux.T != MUX {
		return nil
	}

	if len(args) != 0 {
		return nil
	}

	return &instance{
		T:          ROUTER,
		superRoute: nil,
		pathPrefix: "",
		methods:    nil,
	}
}

// HandleFunc emulate behavior of mux.Router.HandleFunc
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.HandleFunc
func (router *instance) HandleFunc(args ...ast.Expr) *instance {
	if router.T != ROUTER {
		return nil
	}

	if len(args) != 2 {
		return nil
	}

	lit, ok := args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}

	path, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil
	}

	return &instance{
		T:          ROUTE,
		superRoute: router.superRoute,
		pathPrefix: router.pathPrefix,
		methods:    copyMethods(router.methods),
		path:       path,
		isHandled:  true,
	}
}

// PathPrefix emulate behavior of mux.Router.PathPrefix and mux.Route.PathPrefix
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.PathPrefix
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.PathPrefix
func (inst *instance) PathPrefix(args ...ast.Expr) *instance {
	if inst.T != ROUTER && inst.T != ROUTE {
		return nil
	}

	if len(args) != 1 {
		return nil
	}

	lit, ok := args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}

	pathPrefix, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil
	}

	switch inst.T {
	case ROUTER:
		return &instance{
			T:          ROUTE,
			superRoute: inst.superRoute,
			pathPrefix: pathPrefix,
			methods:    copyMethods(inst.methods),
			path:       inst.path,
			isHandled:  false,
		}

	case ROUTE:
		inst.pathPrefix = pathPrefix
		return inst
	}

	return nil
}

// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.Path
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.Path
func (inst *instance) Path(args ...ast.Expr) *instance {
	if inst.T != ROUTER && inst.T != ROUTE {
		return nil
	}

	if len(args) != 1 {
		return nil
	}

	lit, ok := args[0].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return nil
	}

	path, err := strconv.Unquote(lit.Value)
	if err != nil {
		return nil
	}

	switch inst.T {
	case ROUTER:
		return &instance{
			T:          ROUTE,
			superRoute: inst.superRoute,
			pathPrefix: inst.pathPrefix,
			methods:    copyMethods(inst.methods),
			path:       path,
			isHandled:  false,
		}

	case ROUTE:
		inst.pathPrefix = path
		return inst
	}

	return nil
}

// Methods emulate behavior of mux.Router.Methods and mux.Route.Methods
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Router.Methods
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.Methods
func (inst *instance) Methods(args ...ast.Expr) *instance {
	if inst.T != ROUTER && inst.T != ROUTE {
		return nil
	}

	methods := copyMethods(inst.methods)
	for _, arg := range args {
		lit, ok := arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return nil
		}

		method, err := strconv.Unquote(lit.Value)
		if err != nil {
			return nil
		}

		methods = append(methods, method)
	}

	switch inst.T {
	case ROUTER:
		return &instance{
			T:          ROUTE,
			superRoute: inst.superRoute,
			pathPrefix: inst.pathPrefix,
			methods:    methods,
			path:       "",
			isHandled:  false,
		}

	case ROUTE:
		inst.methods = methods
		return inst
	}

	return nil
}

// HandlerFunc emulate behavior of mux.Route.HandlerFunc
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.HandlerFunc
func (route *instance) HandlerFunc(args ...ast.Expr) *instance {
	if route.T != ROUTE {
		return nil
	}

	if len(args) != 1 {
		return nil
	}

	route.isHandled = true
	return route
}

// Subrouter emulate behavior of mux.Route.Subrouter
// Ref. https://pkg.go.dev/github.com/gorilla/mux#Route.Subrouter
func (route *instance) Subrouter(args ...ast.Expr) *instance {
	if route.T != ROUTE {
		return nil
	}

	if len(args) != 0 {
		return nil
	}

	return &instance{
		T:          ROUTER,
		superRoute: route,
		pathPrefix: "",
		methods:    copyMethods(route.methods),
	}
}

func (inst *instance) GetPath() string {
	buf := inst.getPath()
	buf.WriteString(strings.TrimRight(inst.path, "/"))
	path := buf.String()
	if path == "" {
		return "/"
	}

	return path
}

func (inst *instance) getPath() *bytes.Buffer {
	if inst.superRoute == nil {
		return bytes.NewBuffer([]byte(strings.TrimRight(inst.pathPrefix, "/")))
	}

	buf := inst.superRoute.getPath()
	buf.WriteString(strings.TrimRight(inst.pathPrefix, "/"))
	return buf
}

func copyMethods(methods []string) []string {
	if methods == nil {
		return nil
	}

	return methods[:]
}
