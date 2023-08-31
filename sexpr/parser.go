package sexpr

import (
	"strconv"
	"strings"
	"unicode"
)

const whitespace = " \n\t"

func parseInt(s string) (value Expr, remains string, ok bool) {
	sign := 1
	signStr, s, _ := oneOf(matchString("+"), matchString("-"))(s)
	if signStr == "-" {
		sign = -1
	}
	var accum string
	for _, c := range []rune(s) {
		if !unicode.IsDigit(c) {
			break
		}
		accum = accum + string(c)
	}
	if accum == "" {
		return 0, s, false
	}

	n, err := strconv.Atoi(accum)
	if err != nil {
		return 0, s, false
	}
	return sign * n, s[len(accum):], true
}

func parseSymbol(s string) (value Expr, remains string, ok bool) {
	const allowedSymbolChars = "+-*=?/abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var accum string
	for _, c := range []rune(s) {
		if !strings.ContainsRune(allowedSymbolChars, c) {
			break
		}
		ok = true
		accum = accum + string(c)
	}
	return Symbol(accum), s[len(accum):], ok
}

func parseBool(s string) (value Expr, remains string, ok bool) {
	value, remains, ok = oneOf(
		matchString("#t"),
		matchString("#f"),
	)(s)
	return value == "#t", remains, ok
}

func parseString(s string) (value Expr, remains string, ok bool) {
	var accum string
	s, ok = skipRune(s, `"`)
	if !ok {
		return
	}
	for _, c := range []rune(s) {
		if c == '"' {
			break
		}
		ok = true
		accum = accum + string(c)
	}
	remains = s[1+len(accum):]
	return accum, remains, ok
}

// use sepBy
func parseList(s string) (value Expr, remains string, ok bool) {
	remains, ok = skipRune(s, "(")
	if !ok {
		return
	}

	var list []Expr
	for {
		s, _ = skipManyRune(s, whitespace)
		value, remains, ok = Parse(remains)
		if !ok {
			break
		}
		list = append(list, value)

		// optional space between elements (not for last element)
		// FIXME: many spaces
		remains, ok = skipManyRune(remains, whitespace)
		if !ok {
			break
		}
	}

	remains, ok = skipRune(remains, ")")
	if !ok {
		return
	}
	return list, remains, ok
}

func parseQuotedExpr(s string) (value Expr, remains string, ok bool) {
	s, ok = skipRune(s, `'`)
	if !ok {
		return
	}
	innerExpr, remains, ok := Parse(s)
	if !ok {
		return nil, s, false
	}
	return List(Symbol("quote"), innerExpr), remains, true
}

func Parse(s string) (value Expr, remains string, ok bool) {
	s, _ = skipManyRune(s, whitespace)
	return oneOf(parseInt, parseQuotedExpr, parseSymbol, parseBool, parseString, parseList)(s)
}

func MustParse(s string) Expr {
	value, _, ok := Parse(s)
	if !ok {
		panic("parsing failed: " + s)
	}
	return value
}
