package quden

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"strconv"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/framework"
	"github.com/progfay/quden/visitor"
)

func formatBundle(regexp, name string) string {
	return fmt.Sprintf("[[bundle]]\nregexp = %q\nname = %q", regexp, name)
}

var frameworks = []*framework.Framework{
	echo.New(),
}

func findMatchFramework(path string) *framework.Framework {
	for _, framework := range frameworks {
				if framework.MatchImportPath(path) {
					return framework
				}
			}
			return nil
}

func Run(w io.Writer, files []string) {
	v := visitor.New()

	for _, file := range files {
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, file, nil, parser.Mode(0))

		for _, importSpec := range f.Imports {
			path, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				log.Printf("%v\n", err)
				continue
			}

			framework := findMatchFramework(path)
			if framework != nil {
				v.AddNodeConverter(framework.NewNodeConverter())
			}
		}

		ast.Walk(v, f)
	}

	endpoints := v.GetEndpoints()
	for _, endpoint := range endpoints {
		fmt.Fprintln(w, endpoint.String())
	}
}
