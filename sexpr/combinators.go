package sexpr

import "strings"

type Parser = func(s string) (value Expr, remains string, ok bool)

func matchString(toMatch string) Parser {
	return func(s string) (value Expr, remains string, ok bool) {
		if strings.HasPrefix(s, toMatch) {
			return toMatch, s[len(toMatch):], true
		}
		return "", s, false
	}
}

func skipRune(s string, runesToSkip string) (remains string, ok bool) {
	if len(s) == 0 {
		return s, false
	}
	firstRune := []rune(s)[0]
	if strings.ContainsRune(runesToSkip, firstRune) {
		return string([]rune(s)[1:]), true
	}
	return s, false
}

func skipManyRune(s string, toSkip string) (remains string, ok bool) {
	ok = true
	for ok {
		s, ok = skipRune(s, toSkip)
	}
	return s, true
}

func oneOf(parsers ...Parser) Parser {
	return func(s string) (value Expr, remains string, ok bool) {
		for _, parser := range parsers {
			value, remains, ok = parser(s)
			if ok {
				return value, remains, ok
			}
		}
		return value, s, false
	}
}
