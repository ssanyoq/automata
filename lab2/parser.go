package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ssanyoq/automata-uni/lab2/util"
)

type Parser struct {
	// wow so DI so mindful
	lexer *Lexer

	// stacks for stacking
	fragmentsStack *util.Stack[*Automata]
	operatorsStack *util.Stack[Token]

	// To chill for a little bit and only return errors in main-ish functions
	// also I saw this approach in golang parser, so yea
	Errors []error

	// insideParenths int
	// nextCaptGroup  int
}

func combineErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var out strings.Builder
	out.WriteString("following errors occured while parser was working: ")
	for i, v := range errs {
		if i != 0 {
			out.WriteString("; ")
		}
		out.WriteString(v.Error())
	}
	return errors.New(out.String())
}

func NewParser(l *Lexer) *Parser {
	return &Parser{
		lexer:          l,
		fragmentsStack: util.NewStack[*Automata](),
		operatorsStack: util.NewStack[Token](),

		Errors: []error{},
	}
}

func (p *Parser) popOp() {
	op, ok := p.operatorsStack.Pop()
	if !ok {
		p.Errors = append(p.Errors, errors.New("wanted to pop an operator, but couldn't recieve one"))
		return
	}
	args, err := OpRequiresArgs(op)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return
	}
	switch args {
	case 2:
		right, ok := p.fragmentsStack.Pop()
		if !ok {
			p.Errors = append(p.Errors, errors.New("expected operand"))
			return
		}
		left, ok := p.fragmentsStack.Pop()
		if !ok {
			p.Errors = append(p.Errors, errors.New("expected operand"))
			return
		}
		switch op {
		case Concat:
			p.fragmentsStack.Push(ConcatAutomata(left, right))
		case Or:
			p.fragmentsStack.Push(OrAutomata(left, right))
		default:
			panic("it's probably this prognostic guy isn't it")
		}
	case 1:
		operand, ok := p.fragmentsStack.Pop()
		if !ok {
			p.Errors = append(p.Errors, errors.New("expected operand"))
			return
		}
		switch op {
		case Kleene:
			p.fragmentsStack.Push(KleeneeAutomata(operand))
		case PositiveClosure:
			p.fragmentsStack.Push(PositiveClosureAutomata(operand))
		default:
			panic("weird unary operator")
		}
	default:
		p.Errors = append(p.Errors, fmt.Errorf("operation with %d operands is not supported", args))
		return
	}
}

// Called when either the EOL or closed parenthesis is found
func (p *Parser) popTillStop() {
	op, ok := p.operatorsStack.Peek()
	for op != OpenParenthesis && ok {
		p.popOp()
		if len(p.Errors) != 0 {
			return
		}
		op, ok = p.operatorsStack.Peek()
	}
	if op == OpenParenthesis {
		p.operatorsStack.Pop() // to remove them
		// p.fragmentsStack.Push(&CaptureGroupNode{Number: -1, Child: frag}) // do something like this in automatas
	}

}

// Called when parser encounters '{' to create RangeRepeat fragment
func (p *Parser) parseBraces() {
	tok, fromVal := p.lexer.Next()
	var fromStr strings.Builder
	for tok == Digit {
		fromStr.WriteRune(fromVal)
		tok, fromVal = p.lexer.Next()
	}

	if tok != Comma {
		p.Errors = append(p.Errors, fmt.Errorf("at rr expected ',', but got: '%c'", fromVal))
		return
	}
	p.lexer.SkipWhitespaces()

	tok, toVal := p.lexer.Next()
	var toStr strings.Builder
	for tok == Digit {
		toStr.WriteRune(toVal)
		tok, toVal = p.lexer.Next()
	}

	if tok != ClosedBrace {
		p.Errors = append(p.Errors, fmt.Errorf("at rr expected '}', but got: '%c'", fromVal))
		return
	}

	var from, to int
	var err error
	if fromStr.Len() != 0 {
		from, err = strconv.Atoi(fromStr.String())
		if err != nil {
			p.Errors = append(p.Errors, err)
			return
		}
	} else {
		from = 0
	}
	if toStr.Len() != 0 {
		to, err = strconv.Atoi(toStr.String())
		if err != nil {
			p.Errors = append(p.Errors, err)
			return
		}
	} else {
		to = -1
	}
	frag, ok := p.fragmentsStack.Pop()
	if !ok {
		p.Errors = append(p.Errors, errors.New("operator {} needs one operand, but none were found"))
		return
	}
	p.fragmentsStack.Push(RepeatAutomata(frag, from, to))
}

// Called when parser encounters '[' to create RangeRepeat fragment
func (p *Parser) parseBrackets() {
	from, fromVal := p.lexer.Next()
	if from == ClosedBracket {
		p.Errors = append(p.Errors, fmt.Errorf("at [] expected characters, got: %c", fromVal))
		return
	}
	dash, _ := p.lexer.Next()
	if dash != Minus {
		p.Errors = append(p.Errors, fmt.Errorf("at [] expected -, got: %c", fromVal))
		return
	}
	to, toVal := p.lexer.Next()
	if to == ClosedBracket {
		p.Errors = append(p.Errors, fmt.Errorf("at [] expected non-special characters, got: %c", toVal))
		return
	}
	if fromVal > toVal {
		p.Errors = append(p.Errors, fmt.Errorf("'%c' is greater than '%c'", fromVal, toVal))
		return
	}
	p.lexer.Next()
	// p.fragmentsStack.Push(&CharacterRangeNode{
	// 	From: fromVal,
	// 	To:   toVal,
	// })
	p.fragmentsStack.Push(RangeCharAutomata(fromVal, toVal))
}

func (p *Parser) pushPopPriority(nextOp Token) {
	defer p.operatorsStack.Push(nextOp)

	nextPriority, err := OpPriority(nextOp)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return
	}
	headOp, ok := p.operatorsStack.Peek()
	if !ok { // empty stack
		return
	}
	if headOp == OpenParenthesis {
		return
	}
	headPriority, err := OpPriority(headOp)
	if err != nil {
		p.Errors = append(p.Errors, err)
		return
	}
	if headPriority > nextPriority {
		p.popTillStop()
	}
}

func (p *Parser) BuildNFA() (*Automata, error) {
	tok, sym := p.lexer.Next()
	openParenths := 0

	insertConcat := true
	for tok != EOS && len(p.Errors) == 0 {
		switch tok {
		case Character, Comma, Digit, Whitespace: // normies
			if insertConcat && p.fragmentsStack.Size() != 0 {
				p.pushPopPriority(Concat) // insert implicit concat
			}
			p.fragmentsStack.Push(CharAutomata(sym))
		case OpenParenthesis:
			if insertConcat && p.fragmentsStack.Size() != 0 {
				p.pushPopPriority(Concat) // insert implicit concat
			}
			openParenths++
			p.operatorsStack.Push(OpenParenthesis)
		case OpenBrace: // time for using that a{1,2} operator
			p.parseBraces()
		case OpenBracket:
			p.parseBrackets()
		case Prognostic, Concat, Or, PositiveClosure, Kleene:
			p.pushPopPriority(tok)
		case ClosedParenthesis:
			if openParenths < 0 {
				return nil, errors.New("unmatched ')' found")
			}
			p.popTillStop()
			if len(p.Errors) != 0 {
				return nil, combineErrors(p.Errors)
			}
			openParenths--
		}
		args, _ := OpRequiresArgs(tok)
		insertConcat = (args != 2 && tok != OpenParenthesis)
		tok, sym = p.lexer.Next()
	}
	p.popTillStop()
	out, ok := p.fragmentsStack.Pop()
	if !ok {
		return nil, errors.New("no fragments to return")
	}
	return out, combineErrors(p.Errors)
}
