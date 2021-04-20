package quden

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"strconv"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/framework"
	"github.com/progfay/quden/goji"
)

func formatBundle(regexp, name string) string {
	return fmt.Sprintf("[[bundle]]\nregexp = %q\nname = %q", regexp, name)
}

var frameworks = []framework.Framework{
	echo.New(),
	goji.New(),
}

func findMatchFramework(path string) framework.Framework {
	for _, framework := range frameworks {
		if framework.MatchImportPath(path) {
			return framework
		}
	}
	return nil
}

func Run(w io.Writer, files []string) {
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
				endpoints := framework.Extract(f)
				for _, endpoint := range endpoints {
					fmt.Fprintln(w, endpoint.String())
				}
			}
		}
	}
}
