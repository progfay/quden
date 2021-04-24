package quden

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"strconv"

	"github.com/progfay/quden/echo"
	"github.com/progfay/quden/goji"
	"github.com/progfay/quden/gorilla"
	"github.com/progfay/quden/util"
)

func formatBundle(regexp, name string) string {
	return fmt.Sprintf("[[bundle]]\nregexp = %q\nname = %q", regexp, name)
}

var utils = []util.Framework{
	echo.New(),
	goji.New(),
	gorilla.New(),
}

func findMatchFramework(path string) util.Framework {
	for _, util := range utils {
		if util.MatchImportPath(path) {
			return util
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

			util := findMatchFramework(path)
			if util != nil {
				endpoints := util.Extract(f)
				for _, endpoint := range endpoints {
					fmt.Fprintln(w, endpoint.String())
				}
			}
		}
	}
}
