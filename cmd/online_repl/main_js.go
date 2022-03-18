package main

import (
	"syscall/js"

	"github.com/adzeitor/goscheme/scheme"
	"github.com/adzeitor/goscheme/sexpr"
)

func main() {
	js.Global().Set("schemeEval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := scheme.Eval(args[0].String())
		output := sexpr.Print(result)
		return output
	}))
	select {} // Code must not finish
}
