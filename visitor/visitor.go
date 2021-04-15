package visitor

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type visitor struct{
	w io.Writer
}

func isHTTPMethod(str string) bool {
	methodSet := map[string]struct{}{
		http.MethodGet:     {},
		http.MethodHead:    {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodDelete:  {},
		http.MethodConnect: {},
		http.MethodOptions: {},
		http.MethodTrace:   {},
	}
	_, ok := methodSet[str]
	return ok
}

func New(w io.Writer) *visitor {
	return &visitor{w: w}
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return v
	}

	if len(callExpr.Args) < 1 {
		return v
	}

	firstArg, ok := callExpr.Args[0].(*ast.BasicLit)
	if !ok || firstArg.Kind != token.STRING {
		return v
	}
	path, err := strconv.Unquote(firstArg.Value)
	if err != nil || !strings.HasPrefix(path, "/") {
		return v
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return v
	}

	name := strings.ToUpper(selectorExpr.Sel.Name)
	if !isHTTPMethod(name) {
		return v
	}

	fmt.Fprintf(v.w, "%s %s\n", name, path)
	return v
}
