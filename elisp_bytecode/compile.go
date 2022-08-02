package elisp_bytecode

import (
	"github.com/adzeitor/goscheme/sexpr"
)

const (
	StackRef1  = 1
	StackRef2  = 2
	StackRef3  = 3
	Call0      = 32
	Call1      = 33
	Call2      = 34
	Call7      = 39
	Add1       = 94
	Plus       = 92
	Mult       = 95
	Return     = 135
	Dup        = 137
	Constant0  = 192
	Constant1  = 193
	Constant2  = 194
	Constant3  = 195
	Constant4  = 196
	Constant64 = 255
)

func Exec(prog []byte, constants []sexpr.Expr) sexpr.Expr {
	// TODO: add program counter (PC)?
	stack := make([]sexpr.Expr, 0, 10)
	for _, instr := range prog {
		switch {
		case instr >= Constant0:
			stack = append(stack, constants[instr-Constant0])
		//case instr == Constant0 && nextInstr == Constant1:
		//	stack = append(stack, constants[0])
		//	stack = append(stack, constants[1])
		//case instr == Constant1:
		//	stack = append(stack, constants[1])
		//case instr == Constant2:
		//	stack = append(stack, constants[2])
		case instr == StackRef1:
			stack = append(stack, stack[0])
		case instr == StackRef2:
			stack = append(stack, stack[1])
		case instr == StackRef3:
			stack = append(stack, stack[2])
		case instr == Plus:
			result := stack[0].(int) + stack[1].(int)
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		case instr == Mult:
			result := stack[0].(int) * stack[1].(int)
			stack = stack[:len(stack)-2]
			stack = append(stack, result)
		case instr == Return:
			result := stack[0]
			stack = stack[:len(stack)-1]
			return result
		}
	}
	return nil
}

func Compile(object sexpr.Expr) ([]byte, []sexpr.Expr) {
	var constants []sexpr.Expr
	compiled := compile(object, &constants)
	return compiled, constants
}

func compile(object sexpr.Expr, constants *[]sexpr.Expr) []byte {
	// TODO: add program counter (PC)?
	var compiled []byte
	switch v := object.(type) {
	case int:
		return makeConstant(v, constants)
	case sexpr.Symbol:
		return makeConstant(v, constants)
	case []sexpr.Expr:
		if v[0] == sexpr.Symbol("+") {
			compiled = append(compiled, compile(v[1], constants)...)
			compiled = append(compiled, compile(v[2], constants)...)
			compiled = append(compiled, Plus)
			// FIXME: Return is not needed
			compiled = append(compiled, Return)
			return compiled
		} else if v[0] == sexpr.Symbol("*") {
			compiled = append(compiled, compile(v[1], constants)...)
			compiled = append(compiled, compile(v[2], constants)...)
			compiled = append(compiled, Mult)
			return compiled
		} else { // function application
			fn := v[0]
			args := v[1:]
			compiled = append(compiled, compile(fn, constants)...)
			for _, arg := range args {
				compiled = append(compiled, compile(arg, constants)...)
			}
			compiled = append(compiled, Call0+byte(len(args)))
			return compiled
		}
	}
	return nil
}

func makeConstant(value sexpr.Expr, constants *[]sexpr.Expr) []byte {
	nextConstant := Constant0 + byte(len(*constants))
	*constants = append(*constants, value)
	return []byte{nextConstant}
}

//func compileInt(value int) ([]byte, sexpr.Expr) {
//
//}
