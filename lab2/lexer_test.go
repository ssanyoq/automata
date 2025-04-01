package main

import (
	"testing"
)

func TestLexerNext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			token Token
			char  rune
		}
	}{
		{
			name:  "single character",
			input: "a",
			expected: []struct {
				token Token
				char  rune
			}{
				{Character, 'a'},
				{EOS, ' '},
			},
		},
		{
			name:  "characters",
			input: "hello",
			expected: []struct {
				token Token
				char  rune
			}{
				{Character, 'h'},
				{Character, 'e'},
				{Character, 'l'},
				{Character, 'l'},
				{Character, 'o'},
				{EOS, ' '},
			},
		},
		{
			name:  "expression",
			input: "ab+c{}[]())%",
			expected: []struct {
				token Token
				char  rune
			}{
				{Character, 'a'},
				{Character, 'b'},
				{PositiveClosure, '+'},
				{Character, 'c'},
				{OpenBrace, '{'},
				{ClosedBrace, '}'},
				{OpenBracket, '['},
				{ClosedBracket, ']'},
				{OpenParenthesis, '('},
				{ClosedParenthesis, ')'},
				{ClosedParenthesis, ')'},
				{EOS, ' '},
			},
		},
		{
			name:  "escaping",
			input: "%a.b%*",
			expected: []struct {
				token Token
				char  rune
			}{
				{Character, 'a'},
				{Concat, '.'},
				{Character, 'b'},
				{Character, '*'},
				{EOS, ' '},
			},
		},
		{
			name:  "multiple escape characters",
			input: "%%%%a",
			expected: []struct {
				token Token
				char  rune
			}{
				{Character, '%'},
				{Character, '%'},
				{Character, 'a'},
				{EOS, ' '},
			},
		},
		{
			name:  "empty",
			input: "",
			expected: []struct {
				token Token
				char  rune
			}{
				{EOS, ' '},
			},
		},
		{
			name:  "space saving",
			input: "a   .%  ",
			expected: []struct {
				token Token
				char  rune
			}{
				{Character, 'a'},
				{Whitespace, ' '},
				{Whitespace, ' '},
				{Whitespace, ' '},
				{Concat, '.'},
				{Character, ' '},
				{Whitespace, ' '},
				{EOS, ' '},
			},
		},

		{
			name:  "semi-special characters",
			input: "1,  2%3-",
			expected: []struct {
				token Token
				char  rune
			}{
				{Digit, '1'},
				{Comma, ','},
				{Whitespace, ' '},
				{Whitespace, ' '},
				{Digit, '2'},
				{Character, '3'},
				{Minus, '-'},
				{EOS, ' '},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			for i, expected := range tt.expected {
				token, char := lexer.Next()
				if token != expected.token {
					t.Errorf("Test %s: Expected token at index %d to be %v, got %v", tt.name, i, expected.token, token)
				}
				if char != expected.char {
					t.Errorf("Test %s: Expected char at index %d to be %c, got %c", tt.name, i, expected.char, char)
				}
			}
		})
	}
}

func TestSkipWhitespaces(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		start   int
		wantPos int
	}{
		{
			name:    "No leading whitespace",
			input:   "abc",
			start:   0,
			wantPos: 0,
		},
		{
			name:    "Single leading whitespace",
			input:   " abc",
			start:   0,
			wantPos: 1,
		},
		{
			name:    "Multiple leading whitespaces",
			input:   "    abc",
			start:   0,
			wantPos: 4,
		},
		{
			name:    "Only whitespaces",
			input:   "    ",
			start:   0,
			wantPos: 4,
		},
		{
			name:    "Empty string",
			input:   "",
			start:   0,
			wantPos: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			lexer.pos = tt.start
			lexer.SkipWhitespaces()
			if lexer.pos != tt.wantPos {
				t.Errorf("Lexer position after SkipWhitespaces = %d, want %d", lexer.pos, tt.wantPos)
			}
		})
	}
}
