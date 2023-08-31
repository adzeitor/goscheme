package scheme

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/adzeitor/goscheme/sexpr"
)

func TestEval(t *testing.T) {
	t.Run("eval int value", func(t *testing.T) {
		assert.Equal(t, 42, Eval(`42`))
		assert.Equal(t, -13, Eval(`-13`))
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
		assert.Equal(t, false, Eval(`(= #t #f)`))
		assert.Equal(t, true, Eval(`(= #t #t)`))
		assert.Equal(t, true, Eval(`(= #f #f)`))
	})

	t.Run("equal", func(t *testing.T) {
		assert.Equal(t, false, Eval(`(= 1 "1")`))
		assert.Equal(t, false, Eval(`(= 'foo "foo")`))
		assert.Equal(t, true, Eval(`(= () ())`))
		assert.Equal(t, true, Eval(`(= '(3 4 5) '(3 4 5))`))
		assert.Equal(t, false, Eval(`(= '(3 4 5) '(3 4 5 6))`))
		assert.Equal(t, false, Eval(`(= () '(3))`))
		assert.Equal(t, false, Eval(`(= '(3) ())`))
	})

	t.Run("if", func(t *testing.T) {
		assert.Equal(t, 5, Eval(`(if (= 3 3) 5 0)`))
		assert.Equal(t, 0, Eval(`(if (= 3 4) 5 0)`))
	})

	t.Run("cond", func(t *testing.T) {
		assert.Equal(t, true, Eval(`(= 5 (cond ((= 3 3) 5) (else 88)))`))
		assert.Equal(t, true, Eval(`(= 88 (cond ((= 3 6) 5) (else 88)))`))
		assert.Equal(t, true, Eval(`(= 88 (cond ((= 3 4) 5) (#f 16) (#t 88)))`))
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
		assert.Equal(t, true, Eval(`(= 1 (car (cons 1 '(2 3))))`))
		assert.Equal(t, true, Eval(`(= 2 (car (cdr (cons 1 '(2 3)))))`))
	})

	t.Run("list?", func(t *testing.T) {
		assert.Equal(t, true, Eval(`(= #t (list? '(4 5 6)))`))
		assert.Equal(t, true, Eval(`(= #t (list? ()))`))
		assert.Equal(t, true, Eval(`(= #f (list? 5))`))
		assert.Equal(t, true, Eval(`(= #f (list? 'foo))`))
		assert.Equal(t, true, Eval(`(= #f (list? "bar"))`))
	})

	t.Run("symbol?", func(t *testing.T) {
		assert.Equal(t, true, Eval(`(= #f (symbol? '(4 5 6)))`))
		assert.Equal(t, true, Eval(`(= #f (symbol? ()))`))
		assert.Equal(t, true, Eval(`(= #f (symbol? 5))`))
		assert.Equal(t, true, Eval(`(= #t (symbol? 'foo))`))
		assert.Equal(t, true, Eval(`(= #t (symbol? (quote foo)))`))
		assert.Equal(t, true, Eval(`(= #f (symbol? "bar"))`))
		assert.Equal(t, true, Eval(`(= #f (symbol? (quote (a b))))`))
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
		assert.Equal(
			t,
			true,
			Eval(`
				(=
					'(+ '(4 5 6))
					(quote (+ (quote (4 5 6)))))
			`))
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

	t.Run("several arguments lambda", func(t *testing.T) {
		assert.Equal(t, 33, Eval(`((lambda (x y z)  (+ z (* x y))) 5 6 3)`))
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

	t.Run("apply should not extend scope", func(t *testing.T) {
		// arrange
		Eval(`(define hack   (lambda (x) (+ x secret)))`)
		Eval(`(define safe   (lambda (secret) (hack 23)))`)

		// act
		result := Eval(`(safe 42)`)

		// assert
		assert.Contains(t, result, "exception: Unbound variable: secret")
	})

	t.Run("original lambda environment should remain after calling function", func(t *testing.T) {
		// act
		result := Eval(`
       ((lambda (y) 
   			(+ ((lambda (x) y) y)
               ((lambda (x) y) y)))
        4)`)

		// assert
		assert.Equal(t, 8, result)
	})

	t.Run("can create useless lambda :D", func(t *testing.T) {
		t.Skip()
		prog := `
			(lambda (arg)
					(if (= arg 'hello)
						'hi
						(lambda (x) unbound)))
		`
		assert.Equal(t, sexpr.Symbol("hi"), Eval(prog))
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

	// https://www.youtube.com/watch?v=OyfBQmvr2Hc
	t.Run("most beautiful program ever (meta circular evaluator)", func(t *testing.T) {
		Eval(`
			(define eval-expr
				(lambda (expr env)
					(cond
						((symbol? expr) (env expr))
						((= 'lambda (car expr))
							(lambda (arg)
							  (eval-expr (car (cdr (cdr expr)))
								(lambda (y)
									(if (= (car (car (cdr expr))) y)
										arg
										(env y))))))
						(else
							((eval-expr (car expr) env)
								(eval-expr (car (cdr expr)) env))))))
		`)
		prog := `
			(eval-expr 
                '((lambda (n) n) hello)
				(lambda (arg1)
					(if (= arg1 'hello)
						'hi
						unbound))))
		`
		assert.Equal(t, sexpr.Symbol("hi"), Eval(prog))
	})

	t.Run("error on unbound variable", func(t *testing.T) {
		assert.Contains(t, Eval(`foo`), "exception: Unbound variable: foo")
	})

	t.Run("error on applying non-lambda", func(t *testing.T) {
		assert.Equal(t, `exception: The object "+" is not applicable.`, Eval(`("+" 1 2)`))
	})

	t.Run("error on applying car/cdr to wrong types", func(t *testing.T) {
		assert.Equal(
			t,
			`exception: The object 1, passed as the first argument to car, is not the correct type.`,
			Eval(`(car 1)`),
		)
		assert.Equal(
			t,
			`exception: The object "foo", passed as the first argument to cdr, is not the correct type.`,
			Eval(`(cdr "foo")`),
		)
	})

	// https://stackoverflow.com/questions/526082/in-scheme-whats-the-point-of-set
}

func assertEval(t *testing.T, prog string) {
	t.Helper()
	if Eval(prog) != true {
		t.Errorf("evaluation should return true")
	}
}
