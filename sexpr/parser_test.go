package sexpr

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		result  Expr
		wantErr bool
	}{
		{
			in:     `5`,
			result: 5,
		},
		{
			in:     `"foo"`,
			result: "foo",
		},
		{
			in:     `""`,
			result: "",
		},
		{
			in:     `foo`,
			result: Symbol("foo"),
		},
		{
			in:     `#t`,
			result: true,
		},
		{
			in:     `#f`,
			result: false,
		},
		{
			in:     `()`,
			result: List(),
		},
		{
			in:     `(1 "foo" foo)`,
			result: List(1, "foo", Symbol("foo")),
		},
		{
			in: `(* (+ 1 1) (+ 20 1))`,
			result: List(
				Symbol("*"),
				List(Symbol("+"), 1, 1),
				List(Symbol("+"), 20, 1),
			),
		},
		{
			name:   `no spaces between elements`,
			in:     "(1\n2\n3)",
			result: List(1, 2, 3),
		},
		{
			name: `no whitespaces`,
			in:   "((())(())()())",
			result: List(
				List(List()), List(List()), List(), List(),
			),
		},
		{
			name:   "with many white spaces",
			in:     "\t (  1    2 \t3  \t )\t ",
			result: List(1, 2, 3),
		},
		{
			name:   "with many white spaces and tabs",
			in:     "\t (  1    2\t \t3  \t )\t ",
			result: List(1, 2, 3),
		},
		{
			name:    "unclosed list",
			in:      "( ( 1 2 3 )",
			wantErr: true,
		},
	}

	for _, tt := range cases {
		name := tt.name
		if name == "" {
			name = tt.in
		}
		t.Run(name, func(t *testing.T) {
			got, _, ok := Parse(tt.in)
			assert(t, tt.wantErr, !ok)
			assert(t, tt.result, got)
		})
	}
}

func assert(t *testing.T, want, got interface{}) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("not equal want=%+v got=%+v", want, got)
	}
}
