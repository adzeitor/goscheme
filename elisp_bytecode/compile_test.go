package elisp_bytecode

import (
	"reflect"
	"testing"

	"github.com/adzeitor/goscheme/scheme"
	"github.com/adzeitor/goscheme/sexpr"
)

func TestCompile(t *testing.T) {
	t.Run("simple expression", func(t *testing.T) {
		assertCompile(t, "(+ 1 (* 3 4))", "(1 3 4)", Constant0, Constant1, Constant2, Mult, Plus, Return)
	})
	t.Run("function application", func(t *testing.T) {
		assertCompile(t, "(foo)", "(foo)", Constant0, Call0)
		assertCompile(t, "(foo 34)", "(foo 34)", Constant0, Constant1, Call1)
		assertCompile(t, "(bar 23 42)", "(bar 23 42)", Constant0, Constant1, Constant2, Call2)
	})
}

func assertCompile(t *testing.T, prog string, wantVector string, wantBytes ...byte) {
	t.Helper()
	gotBytes, gotVector := Compile(sexpr.MustParse(prog))
	if !reflect.DeepEqual(gotBytes, wantBytes) {
		t.Errorf("compiled bytes wants %v, but got %v", wantBytes, gotBytes)
	}
	if !reflect.DeepEqual(sexpr.Print(gotVector), wantVector) {
		t.Errorf("compiled vector wants %v, but got %v", wantVector, sexpr.Print(gotVector))
	}
}

func TestExec(t *testing.T) {
	t.Run("compile and exec", func(t *testing.T) {
		prog, vector := Compile(sexpr.MustParse("(+ 1 2)"))
		result := Exec(prog, vector)
		if result != 3 {
			t.Fatal(result)
		}
	})
}

func BenchmarkCompile(b *testing.B) {
	prog, vector := Compile(sexpr.MustParse("(+ 33 (* 44 55))"))
	for i := 0; i < b.N; i++ {
		Exec(prog, vector)
	}
}

func BenchmarkEval(b *testing.B) {
	prog := sexpr.MustParse("(+ 33 (* 44 55))")
	for i := 0; i < b.N; i++ {
		scheme.EvalSexpr(prog)
	}
}
