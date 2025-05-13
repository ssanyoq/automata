package main

import "unicode/utf8"

type Lexer struct {
	input []rune

	pos int
}

// Converts string into slice of runes(tomato->tomato)
func toRunes(str string) []rune {
	result := []rune{}
	bytes := []byte(str)
	for i := 0; i < len(bytes); {
		r, size := utf8.DecodeRune(bytes[i:])
		result = append(result, r)
		i += size
	}
	return result
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input: toRunes(input),
		pos:   0,
	}
}

func (l *Lexer) SkipWhitespaces() {
	if l.pos >= len(l.input) || GetToken(l.input[l.pos]) != Whitespace {
		return
	}
	l.Next()
	for l.pos < len(l.input) && GetToken(l.input[l.pos]) == Whitespace {
		l.pos++
	}
}

// Returns current token and rune, then
// makes a step to the next character if possible
func (l *Lexer) Next() (Token, rune) {
	if l.pos == len(l.input) {
		return EOS, ' '
	}
	char := l.input[l.pos]
	t := GetToken(char)
	// Managing escape characters
	switch t {
	case Escape:
		l.pos++
		if l.pos == len(l.input) {
			return EOS, ' '
		}
		t = Character
		char = l.input[l.pos]
		l.pos++
	default:
		l.pos++
	}
	return t, char
}
