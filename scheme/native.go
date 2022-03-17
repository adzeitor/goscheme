package scheme

import "github.com/adzeitor/goscheme/sexpr"

func addBultin(env Environment) {
	env["+"] = Builtin(plusBuiltin)
	env["-"] = Builtin(minusBuiltin)
	env["*"] = Builtin(multBuiltin)
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
