package scheme

import (
	"fmt"
	"strings"

	"github.com/adzeitor/goscheme/sexpr"
)

var globalEnvironment = DefaultEnvironment()

type Builtin func(args []sexpr.Expr, env Environment) sexpr.Expr

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
	result = eval(parsed, env)
	return result, env
}

func EvalBuffer(s string, env Environment) (result sexpr.Expr, resultEnv Environment) {
	// FIXME: is exceptions here good idea? Maybe move to repl or exception free code?
	defer func() {
		err := recover()
		if err != nil {
			result = fmt.Sprintf("exception: %v", err)
			// FIXME: good?
			resultEnv = env
		}
	}()

	for {
		if strings.TrimSpace(s) == "" {
			break
		}

		parsed, remains, ok := sexpr.Parse(s)
		if !ok {
			return "parse error", env
		}
		result = eval(parsed, env)
		s = remains
	}
	return result, env
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

func applyLambda(lambda Lambda, arguments []sexpr.Expr) sexpr.Expr {
	closureEnv := lambda.Env.Extend(lambda.MakeArgEnv(arguments))
	return eval(lambda.Body, closureEnv)
}

func evalArguments(args []sexpr.Expr, env Environment) []sexpr.Expr {
	evaledArgs := make([]sexpr.Expr, len(args))
	for i, arg := range args {
		evaledArgs[i] = eval(arg, env)
	}
	return evaledArgs
}

func evalList(list []sexpr.Expr, env Environment) sexpr.Expr {
	if len(list) == 0 {
		return sexpr.List()
	}
	head := list[0]

	switch head {
	case sexpr.Symbol("quote"):
		return list[1]
	case sexpr.Symbol("="):
		value1 := eval(list[1], env)
		value2 := eval(list[2], env)
		return sexpr.Equal(value1, value2)
	case sexpr.Symbol("null?"):
		arg := eval(list[1], env)
		if l, ok := arg.([]sexpr.Expr); ok {
			return len(l) == 0
		}
		return false
	case sexpr.Symbol("if"):
		condition := eval(list[1], env)
		if condition.(bool) {
			return eval(list[2], env)
		} else {
			return eval(list[3], env)
		}
	case sexpr.Symbol("define"):
		name := list[1].(sexpr.Symbol)
		body := list[2]
		value := eval(body, env)
		env.Global[name] = value
		return name
	case sexpr.Symbol("cons"):
		car := eval(list[1], env)
		cdr := eval(list[2], env)
		// FIXME: optmize and allocate once
		l := sexpr.List(car)
		// FIXME: second argument can be pair...
		// For example, (cons 4 5)  => (4 . 5)
		l = append(l.([]sexpr.Expr), cdr.([]sexpr.Expr)...)
		return l
	case sexpr.Symbol("cond"):
		result := evalCond(list[1:], env)
		return result
	case sexpr.Symbol("lambda"):
		return Lambda{
			Env:        env.Copy(),
			Parameters: list[1].([]sexpr.Expr),
			Body:       list[2],
		}
	}

	head = eval(head, env)
	switch head.(type) {
	case Builtin:
		return head.(Builtin)(list[1:], env)
	case Lambda:
		arguments := evalArguments(list[1:], env)
		result := applyLambda(head.(Lambda), arguments)
		return result
	default:
		panic(fmt.Sprintf("The object %v is not applicable.", sexpr.Print(head)))
	}
}

func evalCond(conditions []sexpr.Expr, env Environment) sexpr.Expr {
	for _, clause := range conditions {
		if clause.([]sexpr.Expr)[0] == sexpr.Symbol("else") {
			return eval(clause.([]sexpr.Expr)[1], env)
		}

		match := eval(clause.([]sexpr.Expr)[0], env)
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

func eval(expr sexpr.Expr, env Environment) sexpr.Expr {
	switch value := expr.(type) {
	case int:
		return value
	case string:
		return value
	case bool:
		return value
	case sexpr.Symbol:
		v := env.Local[value]
		if v != nil {
			return v
		}

		v = env.Global[value]
		if v != nil {
			return v
		}
		panic("Unbound variable: " + value)
	case []sexpr.Expr:
		return evalList(value, env)
	}
	return nil
}
