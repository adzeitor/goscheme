package scheme

import (
	"fmt"

	"github.com/adzeitor/goscheme/sexpr"
)

var globalEnvironment = DefaultEnvironment()

type Builtin func(args []sexpr.Expr, env Environment) (sexpr.Expr, Environment)

// Eval evaluate expression in global environment. Be careful with several
// instances, in this case use EvalInEnvironment.
func Eval(s string) sexpr.Expr {
	result, newEnvironment := EvalInEnvironment(s, globalEnvironment)
	globalEnvironment = newEnvironment
	return result
}

func EvalInEnvironment(s string, env Environment) (result sexpr.Expr, resultEnv Environment) {
	// FIXME: is exceptions here good idea? Maybe move to repl or exception free code?
	defer func() {
		err := recover()
		if err != nil {
			result = fmt.Sprintf("exception: %v", err)
			// FIXME: good?
			resultEnv = env
		}
	}()

	parsed, _, ok := sexpr.Parse(s)
	if !ok {
		return "parse error", env
	}
	return eval(parsed, env)
}

type Lambda struct {
	Env        Environment
	Parameters []sexpr.Expr
	Body       sexpr.Expr
}

func (lambda Lambda) MakeArgEnv(arguments []sexpr.Expr) Environment {
	env := make(Environment)
	for i, value := range arguments {
		// FIXME: parametrs out of index
		env[lambda.Parameters[i].(sexpr.Symbol)] = value
	}
	return env
}

func applyLambda(lambda Lambda, arguments []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	closureEnv := env.Extend(lambda.Env)
	closureEnv = closureEnv.Extend(lambda.MakeArgEnv(arguments))
	return eval(lambda.Body, closureEnv)
}

func evalArguments(args []sexpr.Expr, env Environment) ([]sexpr.Expr, Environment) {
	evaledArgs := make([]sexpr.Expr, len(args))
	for i, arg := range args {
		evaledArgs[i], env = eval(arg, env)
	}
	return evaledArgs, env
}

func evalList(list []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	if len(list) == 0 {
		return sexpr.List(), env
	}
	head := list[0]

	// FIXME: add cons?
	switch head {
	case sexpr.Symbol("quote"):
		return list[1], env
	case sexpr.Symbol("="):
		value1, _ := eval(list[1], env)
		value2, _ := eval(list[2], env)
		// FIXME: recursive lists comparison?
		return sexpr.Equal(value1, value2), env
	case sexpr.Symbol("null?"):
		arg, env := eval(list[1], env)
		if l, ok := arg.([]sexpr.Expr); ok {
			return len(l) == 0, env
		}
		return false, env
	case sexpr.Symbol("if"):
		condition, _ := eval(list[1], env)
		if condition.(bool) {
			return eval(list[2], env)
		} else {
			return eval(list[3], env)
		}
	case sexpr.Symbol("define"):
		name := list[1].(sexpr.Symbol)
		value, _ := eval(list[2], env)
		env[name] = value
		return name, env
	case sexpr.Symbol("cons"):
		car, env := eval(list[1], env)
		cdr, env := eval(list[2], env)
		// FIXME: optmize and allocate once
		l := sexpr.List(car)
		l = append(l.([]sexpr.Expr), cdr.([]sexpr.Expr)...)
		return l, env
	case sexpr.Symbol("cdr"):
		arg, env := eval(list[1], env)
		// FIXME: panic and copy?
		return arg.([]sexpr.Expr)[1:], env
	case sexpr.Symbol("car"):
		arg, env := eval(list[1], env)
		// FIXME: panic and copy?
		return arg.([]sexpr.Expr)[0], env
	case sexpr.Symbol("lambda"):
		return Lambda{
			Env:        env.Copy(),
			Parameters: list[1].([]sexpr.Expr),
			Body:       list[2],
		}, env
	}

	head, env = eval(head, env)
	switch head.(type) {
	case Builtin:
		return head.(Builtin)(list[1:], env)
	case Lambda:
		arguments, env := evalArguments(list[1:], env)
		return applyLambda(head.(Lambda), arguments, env)
	}
	return nil, env
}

type Environment map[sexpr.Symbol]sexpr.Expr

func (env Environment) Copy() Environment {
	copied := make(Environment, len(env))
	for k, v := range env {
		copied[k] = v
	}
	return copied
}

func (env Environment) Extend(extension Environment) Environment {
	newEnv := env.Copy()
	// add extension to new environment
	for k, v := range extension {
		newEnv[k] = v
	}
	return newEnv
}

func DefaultEnvironment() Environment {
	env := make(Environment)
	addBultin(env)
	return env
}

func eval(expr sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	switch value := expr.(type) {
	case int:
		return value, env
	case string:
		return value, env
	case sexpr.Symbol:
		// FIXME: what if not defined?
		if v := env[value]; v != nil {
			return v, env
		}
	case []sexpr.Expr:
		return evalList(value, env)
	}
	return nil, env
}
