package scheme

import (
	"fmt"

	"github.com/adzeitor/goscheme/sexpr"
)

func addBultin(env Environment) {
	env.Global["+"] = Builtin(plusBuiltin)
	env.Global["-"] = Builtin(minusBuiltin)
	env.Global["*"] = Builtin(multBuiltin)
	env.Global[">"] = Builtin(greaterBuiltin)
	env.Global["<"] = Builtin(lessBuiltin)
	env.Global["car"] = Builtin(carBuiltin)
	env.Global["cdr"] = Builtin(cdrBuiltin)
	env.Global["list?"] = Builtin(isListBuiltin)
	env.Global["symbol?"] = Builtin(isSymbolBuiltin)
	env.Global["do"] = Builtin(doBuiltin)
	env.Global["set!"] = Builtin(setBuiltin)
}

func plusBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	n1 := eval(args[0], env)
	n2 := eval(args[1], env)
	return n1.(int) + n2.(int)
}

func minusBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	n1 := eval(args[0], env)
	n2 := eval(args[1], env)
	return n1.(int) - n2.(int)
}

func multBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	n1 := eval(args[0], env)
	n2 := eval(args[1], env)
	return n1.(int) * n2.(int)
}

func greaterBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	n1 := eval(args[0], env)
	n2 := eval(args[1], env)
	return n1.(int) > n2.(int)
}

func lessBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	n1 := eval(args[0], env)
	n2 := eval(args[1], env)
	return n1.(int) < n2.(int)
}

func carBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	arg := eval(args[0], env)
	list, ok := arg.([]sexpr.Expr)
	if !ok {
		panic(
			fmt.Sprintf(
				"The object %v, passed as the first argument to car, is not the correct type.",
				sexpr.Print(arg),
			))
	}
	return list[0]
}

func cdrBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	arg := eval(args[0], env)
	list, ok := arg.([]sexpr.Expr)
	if !ok {
		panic(
			fmt.Sprintf(
				"The object %v, passed as the first argument to cdr, is not the correct type.",
				sexpr.Print(arg),
			))
	}
	return list[1:]
}

func isListBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	arg := eval(args[0], env)
	_, ok := arg.([]sexpr.Expr)
	return ok
}

func isSymbolBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	arg := eval(args[0], env)
	_, ok := arg.(sexpr.Symbol)
	return ok
}

func doBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	result := sexpr.Expr(nil)
	for _, arg := range args {
		result = eval(arg, env)
	}
	return result
}

func setBuiltin(args []sexpr.Expr, env Environment) sexpr.Expr {
	name := args[0].(sexpr.Symbol)
	value := eval(args[1], env)
	// global?
	if _, ok := env.Global[name]; ok {
		env.Global[name] = value
		return nil
	}
	// local?
	env.Local[name] = value
	return nil
}

// FIXME: maybe change to WithEvalArguments
// or just introduce defmacro
func AddFuncToEnv(env Environment, name string, f Builtin) {
	env.Global[sexpr.Symbol(name)] = Builtin(func(args []sexpr.Expr, env Environment) sexpr.Expr {
		evaledArgs := make([]sexpr.Expr, 0, len(args))
		for _, arg := range args {
			evaledArgs = append(evaledArgs, eval(arg, env))
		}
		return f(evaledArgs, env)
	})
}
