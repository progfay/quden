package main

import (
	"flag"
	"os"

	"github.com/progfay/quden"
)

func main() {
	flag.Parse()
	args := flag.Args()
	quden.Run(os.Stdout, args)
}
