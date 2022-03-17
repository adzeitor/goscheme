package sexpr

import (
	"fmt"
	"strconv"
	"strings"
)

type Expr interface{}

type Symbol string

func List(l ...Expr) Expr {
	return l
}

func IsComplete(s string) bool {
	bracesBalance := 0
	isString := false
	hasSexpr := false
	for _, c := range []rune(s) {
		if !strings.ContainsRune(whitespace, c) {
			hasSexpr = true
		}
		if c == '"' {
			isString = !isString
		}
		if c == '(' && !isString {
			bracesBalance++
		}
		if c == ')' && !isString {
			bracesBalance--
		}
	}
	return bracesBalance <= 0 && hasSexpr && !isString
}

func Print(e Expr) string {
	switch value := e.(type) {
	case int:
		return strconv.Itoa(value)
	case string:
		return `"` + value + `"`
	case Symbol:
		return string(value)
	case bool:
		if value {
			return "#t"
		} else {
			return "#f"
		}
	case []Expr:
		elements := make([]string, 0, len(value))
		for _, element := range value {
			elements = append(elements, Print(element))
		}
		return "(" + strings.Join(elements, " ") + ")"
	default:
		return fmt.Sprintf("UNKNOWN %t", e)
	}
}

func Equal(one Expr, other Expr) bool {
	switch value := one.(type) {
	case int:
		if _, ok := other.(int); !ok {
			return false
		}
		return value == other.(int)
	case string:
		if _, ok := other.(string); !ok {
			return false
		}
		return value == other.(string)
	case Symbol:
		if _, ok := other.(Symbol); !ok {
			return false
		}
		return value == other.(Symbol)
	case bool:
		if _, ok := other.(bool); !ok {
			return false
		}
		return value == other.(bool)
	case []Expr:
		second, ok := other.([]Expr)
		if !ok {
			return false
		}
		if len(value) != len(second) {
			return false
		}
		for i := range value {
			if !Equal(value[i], second[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
