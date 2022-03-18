package scheme

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adzeitor/goscheme/sexpr"
)

func TestEval(t *testing.T) {
	t.Run("eval int value", func(t *testing.T) {
		got := Eval(`42`)
		assert.Equal(t, 42, got)
	})

	t.Run("eval int builtin function", func(t *testing.T) {
		assert.Equal(t, 42, Eval(`(+ 20 22)`))
		assert.Equal(t, 30, Eval(`(* 5 6)`))
		assert.Equal(t, 10, Eval(`(- 20 10)`))

		t.Run("multiple arguments", func(t *testing.T) {
			t.Skip("currently multiple arguments of + is not supported")
			assert.Equal(t, 52, Eval(`(+ 20 22 10)`))
		})

		t.Run("division", func(t *testing.T) {
			t.Skip("unsupported")
			assert.Equal(t, 5, Eval(`(/ 10 2)`))
		})
	})

	t.Run("complex plus", func(t *testing.T) {
		assert.Equal(t, 42, Eval(`(+ (* 10 2) (+ 2 20))`))
	})

	t.Run("boolean", func(t *testing.T) {
		assert.Equal(t, true, Eval(`(= 3 3)`))
		assert.Equal(t, false, Eval(`(= 3 4)`))
	})

	t.Run("equal", func(t *testing.T) {
		assert.Equal(t, false, Eval(`(= 1 "1")`))
		assert.Equal(t, false, Eval(`(= (quote foo) "foo")`))
		assert.Equal(t, true, Eval(`(= () ())`))
		assert.Equal(t, true, Eval(`(= (quote (3 4 5)) (quote (3 4 5)))`))
		assert.Equal(t, false, Eval(`(= (quote (3 4 5)) (quote (3 4 5 6)))`))
		assert.Equal(t, false, Eval(`(= (quote ()) (quote (3)))`))
		assert.Equal(t, false, Eval(`(= (quote (3)) (quote ()))`))
	})

	t.Run("if", func(t *testing.T) {
		assert.Equal(t, 5, Eval(`(if (= 3 3) 5 0)`))
		assert.Equal(t, 0, Eval(`(if (= 3 4) 5 0)`))
	})

	t.Run("null?", func(t *testing.T) {
		assert.Equal(t, true, Eval(`(null? ())`))
		assert.Equal(t, false, Eval(`(null? (quote (1)))`))
		assert.Equal(t, false, Eval(`(null? 5)`))
	})

	t.Run("car cdr cons", func(t *testing.T) {
		assert.Equal(t, 4, Eval(`(car (quote (4 5 6)))`))
		assert.Equal(t,
			sexpr.List(2, 3),
			Eval(`(cdr (quote (1 2 3)))`),
		)
		assert.Equal(t, true, Eval(`(= 1 (car (cons 1 (quote (2 3)))))`))
		assert.Equal(t, true, Eval(`(= 2 (car (cdr (cons 1 (quote (2 3))))))`))
	})

	t.Run("quotes", func(t *testing.T) {
		assert.Equal(
			t,
			sexpr.List(1, 2, 3),
			Eval(`(quote (1 2 3))`),
		)
		assert.Equal(
			t,
			5,
			Eval(`(quote 5)`),
		)
	})

	t.Run("with variables", func(t *testing.T) {
		// arrange
		Eval(`(define forty-two (* 21 2))`)

		// act
		result := Eval(`(+ forty-two 1)`)

		// assert
		assert.Equal(t, 43, result)
	})

	t.Run("simple lambda", func(t *testing.T) {
		assert.Equal(t, 25, Eval(`((lambda (x) (* x x)) 5)`))
	})

	t.Run("define lambda", func(t *testing.T) {
		// arrange
		Eval(`(define square (lambda (x) (* x x)))`)

		// act
		result := Eval(`(square 6)`)

		// assert
		assert.Equal(t, 36, result)
	})

	t.Run("should evaluate arguments of lambda", func(t *testing.T) {
		// arrange
		Eval(`(define square (lambda (x) (* x x)))`)

		// act
		result := Eval(`(square (+ 3 3))`)

		// assert
		assert.Equal(t, 36, result)
	})

	t.Run("pass lambda in lambda lambda", func(t *testing.T) {
		// arrange
		Eval(`(define double (lambda (x) (+ x x)))`)
		Eval(`(define twice  (lambda (fn) (lambda (x) (fn (fn x)))))`)

		// act
		result := Eval(`((twice double) 10)`)

		// assert
		assert.Equal(t, 40, result)
	})

	t.Run("recursion", func(t *testing.T) {
		// arrange
		Eval(`
			(define fact
				(lambda (n)
					(if (= n 0)
						1
						(* n (fact (- n 1))))))
		`)

		// assert
		assert.Equal(t, 120, Eval(`(fact 5)`))
	})

	t.Run("error on unbound variable", func(t *testing.T) {
		assert.Equal(t, "exception: Unbound variable: foo", Eval(`foo`))
	})
}
