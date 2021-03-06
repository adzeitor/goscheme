package scheme

import (
	"fmt"

	"github.com/adzeitor/goscheme/sexpr"
)

func addBultin(env Environment) {
	env.Global["+"] = Builtin(plusBuiltin)
	env.Global["-"] = Builtin(minusBuiltin)
	env.Global["*"] = Builtin(multBuiltin)
	env.Global["car"] = Builtin(carBuiltin)
	env.Global["cdr"] = Builtin(cdrBuiltin)
	env.Global["list?"] = Builtin(isListBuiltin)
	env.Global["symbol?"] = Builtin(isSymbolBuiltin)
}

func plusBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	n1, _ := eval(args[0], env)
	n2, _ := eval(args[1], env)
	return n1.(int) + n2.(int), env
}

func minusBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	n1, _ := eval(args[0], env)
	n2, _ := eval(args[1], env)
	return n1.(int) - n2.(int), env
}

func multBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	n1, _ := eval(args[0], env)
	n2, _ := eval(args[1], env)
	return n1.(int) * n2.(int), env
}

func carBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	arg, env := eval(args[0], env)
	list, ok := arg.([]sexpr.Expr)
	if !ok {
		panic(
			fmt.Sprintf(
				"The object %v, passed as the first argument to car, is not the correct type.",
				sexpr.Print(arg),
			))
	}
	return list[0], env
}

func cdrBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	arg, env := eval(args[0], env)
	list, ok := arg.([]sexpr.Expr)
	if !ok {
		panic(
			fmt.Sprintf(
				"The object %v, passed as the first argument to cdr, is not the correct type.",
				sexpr.Print(arg),
			))
	}
	return list[1:], env
}

func isListBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	arg, env := eval(args[0], env)
	_, ok := arg.([]sexpr.Expr)
	return ok, env
}

func isSymbolBuiltin(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	arg, env := eval(args[0], env)
	_, ok := arg.(sexpr.Symbol)
	return ok, env
}
