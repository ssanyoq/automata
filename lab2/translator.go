package main

import (
	"strings"
)

// checks if given symbol is special. If so, escapes it
func escapeSymbol(sym rune) string {
	var out strings.Builder
	switch sym {
	case '.', '(', ')', '{', '}', '[', ']', '|', '%', '/',
		'*', '+':
		out.WriteByte('%')
	}
	out.WriteRune(sym)
	return out.String()
}
