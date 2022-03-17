package main

import (
	"fmt"
	"syscall/js"

	"github.com/adzeitor/goscheme/scheme"
)

func main() {
	js.Global().Set("schemeEval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := scheme.Eval(args[0].String())
		output := fmt.Sprint(result)
		return output
	}))
	select {} // Code must not finish
}
