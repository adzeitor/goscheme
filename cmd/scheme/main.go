package main

import (
	"os"

	"github.com/adzeitor/goscheme/scheme"
)

func main() {
	scheme.RunRepl(os.Stdin, os.Stdout)
}
