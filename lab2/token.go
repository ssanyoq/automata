package main

import "errors"

type Token int

const (
	EOS Token = iota

	Concat          // '.' or ''
	Prognostic      // '/'
	Or              // '|'
	Kleene          // '*'
	PositiveClosure // '+'

	OpenParenthesis   // '('
	ClosedParenthesis // ')'
	OpenBracket       // '['
	ClosedBracket     // ']'
	OpenBrace         // '{'
	ClosedBrace       // '}'

	Escape // %

	Character  // any non-special (or escaped special) character except ones mentioned below
	Digit      // digits
	Minus      // -
	Comma      // ,
	Whitespace // ' '

)

func GetToken(character rune) Token {
	switch character {
	case '.':
		return Concat
	case '/':
		return Prognostic
	case '|':
		return Or
	case '*':
		return Kleene
	case '+':
		return PositiveClosure
	case '(':
		return OpenParenthesis
	case ')':
		return ClosedParenthesis
	case '[':
		return OpenBracket
	case ']':
		return ClosedBracket
	case '{':
		return OpenBrace
	case '}':
		return ClosedBrace
	case '%':
		return Escape
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return Digit
	case '-':
		return Minus
	case ',':
		return Comma
	case ' ':
		return Whitespace
	default:
		return Character
	}
}

func OpRequiresArgs(t Token) (int, error) {
	switch t {
	case Concat:
		return 2, nil
	case Prognostic:
		return 2, nil
	case Or:
		return 2, nil
	case Kleene:
		return 1, nil
	case PositiveClosure:
		return 1, nil
	}
	return 0, errors.New("unknown operation")
}

// Returns operators priority. Highest priority is 0
func OpPriority(t Token) (int, error) {
	switch t {
	case Concat:
		return 1, nil
	case Prognostic:
		return 2, nil
	case Or:
		return 2, nil
	case Kleene:
		return 0, nil
	case PositiveClosure:
		return 0, nil
	}
	return -1, errors.New("unknown operation")
}
