package quden

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/progfay/quden/visitor"
)

func formatBundle(regexp, name string) string {
	return fmt.Sprintf("[[bundle]]\nregexp = %q\nname = %q", regexp, name)
}

func Run(w io.Writer, files []string) {
	v := visitor.New()

	for _, file := range files {
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, file, nil, parser.Mode(0))
		ast.Walk(v, f)
	}

	endpoints := v.GetEndpoints()
	for _, endpoint := range endpoints {
		fmt.Fprintln(w, endpoint)
	}
}
