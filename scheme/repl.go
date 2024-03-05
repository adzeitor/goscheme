package scheme

import (
	"bufio"
	"fmt"
	"io"

	"github.com/adzeitor/goscheme/sexpr"
)

func scan(scanner *bufio.Scanner) bool {
	return scanner.Scan()
}

func RunRepl(
	env Environment,
	input io.Reader,
	output io.Writer,
) {
	// FIXME: it is not fully suitable here because newline may break sexpr which
	// results in parser error.
	buf := bufio.NewScanner(input)

	var result sexpr.Expr
	line := ""
	fmt.Fprint(output, "> ")
	for scan(buf) {
		line += buf.Text()
		if !sexpr.IsComplete(line) {
			continue
		}
		result, env = EvalInEnvironment(line, env)
		fmt.Fprintln(output, sexpr.Print(result))
		fmt.Fprintln(output)
		line = ""
		fmt.Fprint(output, "> ")
	}
}
