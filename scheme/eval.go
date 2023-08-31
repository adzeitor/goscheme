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
	env := EmptyEnvironment()
	for i, value := range arguments {
		// FIXME: parameters out of index
		parameter, ok := lambda.Parameters[i].(sexpr.Symbol)
		if !ok {
			panic(fmt.Sprintf("invalid parameter %v", lambda.Parameters[i]))
		}
		env.Local[parameter] = value
	}
	return env
}

func applyLambda(lambda Lambda, arguments []sexpr.Expr) (sexpr.Expr, Environment) {
	closureEnv := lambda.Env.Extend(lambda.MakeArgEnv(arguments))
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

	switch head {
	case sexpr.Symbol("quote"):
		return list[1], env
	case sexpr.Symbol("="):
		value1, _ := eval(list[1], env)
		value2, _ := eval(list[2], env)
		return sexpr.Equal(value1, value2), env
	case sexpr.Symbol("null?"):
		arg, _ := eval(list[1], env)
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
		body := list[2]
		value, _ := eval(body, env)
		env.Global[name] = value
		return name, env
	case sexpr.Symbol("cons"):
		car, _ := eval(list[1], env)
		cdr, _ := eval(list[2], env)
		// FIXME: optmize and allocate once
		l := sexpr.List(car)
		// FIXME: second argument can be pair...
		// For example, (cons 4 5)  => (4 . 5)
		l = append(l.([]sexpr.Expr), cdr.([]sexpr.Expr)...)
		return l, env
	case sexpr.Symbol("cond"):
		result, _ := evalCond(list[1:], env)
		return result, env
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
		arguments, _ := evalArguments(list[1:], env)
		result, _ := applyLambda(head.(Lambda), arguments)
		return result, env
	default:
		panic(fmt.Sprintf("The object %v is not applicable.", sexpr.Print(head)))
	}
}

func evalCond(conditions []sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	for _, clause := range conditions {
		if clause.([]sexpr.Expr)[0] == sexpr.Symbol("else") {
			return eval(clause.([]sexpr.Expr)[1], env)
		}

		match, _ := eval(clause.([]sexpr.Expr)[0], env)
		if match == true {
			return eval(clause.([]sexpr.Expr)[1], env)
		}
	}
	panic("no match in cond")
}

type Environment struct {
	Global map[sexpr.Symbol]sexpr.Expr
	Local  map[sexpr.Symbol]sexpr.Expr
}

func EmptyEnvironment() Environment {
	return Environment{
		Global: make(map[sexpr.Symbol]sexpr.Expr),
		Local:  make(map[sexpr.Symbol]sexpr.Expr),
	}
}

func (env Environment) Copy() Environment {
	copiedLocal := make(map[sexpr.Symbol]sexpr.Expr, len(env.Local))
	for k, v := range env.Local {
		copiedLocal[k] = v
	}
	return Environment{
		Global: env.Global,
		Local:  copiedLocal,
	}
}

func (env Environment) Extend(extension Environment) Environment {
	newEnv := env.Copy()
	// add extension to new environment
	for k, v := range extension.Local {
		newEnv.Local[k] = v
	}
	return newEnv
}

func DefaultEnvironment() Environment {
	env := EmptyEnvironment()
	addBultin(env)
	return env
}

func EvalSexpr(expr sexpr.Expr) (sexpr.Expr, Environment) {
	return eval(expr, globalEnvironment)
}

func eval(expr sexpr.Expr, env Environment) (sexpr.Expr, Environment) {
	switch value := expr.(type) {
	case int:
		return value, env
	case string:
		return value, env
	case bool:
		return value, env
	case sexpr.Symbol:
		v := env.Local[value]
		if v != nil {
			return v, env
		}

		v = env.Global[value]
		if v != nil {
			return v, env
		}
		panic("Unbound variable: " + value)
	case []sexpr.Expr:
		return evalList(value, env)
	}
	return nil, env
}
