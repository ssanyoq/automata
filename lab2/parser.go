package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	// wow so DI so mindful
	lexer *Lexer

	// stacks for stacking
	fragmentsStack *Stack[Node]
	operatorsStack *Stack[Token]

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
		fragmentsStack: NewStack[Node](),
		operatorsStack: NewStack[Token](),

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
		p.fragmentsStack.Push(&BinaryOpNode{Left: left, Right: right, Operation: op})
	case 1:
		operand, ok := p.fragmentsStack.Pop()
		if !ok {
			p.Errors = append(p.Errors, errors.New("expected operand"))
			return
		}
		p.fragmentsStack.Push(&UnaryOpNode{Child: operand, Operation: op})
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
		frag, ok := p.fragmentsStack.Pop()
		if !ok {
			p.Errors = append(p.Errors, errors.New("after popping parentheses there are no more nodes left somehow"))
			return
		}
		p.fragmentsStack.Push(&CaptureGroupNode{Number: -1, Child: frag}) // TODO correct numbers
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

	var toStr strings.Builder
	for tok == Digit {
		toStr.WriteRune(fromVal)
		tok, fromVal = p.lexer.Next()
	}
	if tok != ClosedBrace {
		p.Errors = append(p.Errors, fmt.Errorf("at rr expected '}', but got: '%c'", fromVal))
		return
	}
	var from, to *int
	if fromStr.Len() != 0 {
		fromVal, err := strconv.Atoi(fromStr.String())
		from = &fromVal
		if err != nil {
			p.Errors = append(p.Errors, err)
			return
		}
		toVal, err := strconv.Atoi(fromStr.String())
		to = &toVal
		if err != nil {
			p.Errors = append(p.Errors, err)
			return
		}
	}
	frag, ok := p.fragmentsStack.Pop()
	if !ok {
		p.Errors = append(p.Errors, errors.New("operator {} needs one operand, but none were found"))
		return
	}
	p.fragmentsStack.Push(&RangeRepeatNode{
		From:  from,
		To:    to,
		Child: frag,
	})
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
	p.fragmentsStack.Push(&CharacterRangeNode{
		From: fromVal,
		To:   toVal,
	})
}

func (p *Parser) BuildAST() (*Node, error) {
	tok, sym := p.lexer.Next()
	openParenths := 0

	// To insert '.' if wasn't
	// wasBinary := true
	for tok != EOS {
		switch tok {
		case Character:
			p.fragmentsStack.Push(&CharNode{Character: sym})
		case OpenParenthesis:
			openParenths++
			p.operatorsStack.Push(OpenParenthesis)
		case OpenBrace: // time for using that a{1,2} operator
			p.parseBraces()
		case ClosedParenthesis:
			if openParenths < 0 {
				return nil, errors.New("unmatched ')' found")
			}
			p.popTillStop()
			if p.Errors != nil {
				return nil, combineErrors(p.Errors)
			}
			openParenths--
			frag, ok := p.fragmentsStack.Pop()
			if !ok {
				return nil, errors.New("somehow poptillstop left 0 fragments")
			}
			p.fragmentsStack.Push(&CaptureGroupNode{Child: frag, Number: -1}) // TODO use correct numbers
		}
		tok, sym = p.lexer.Next()
	}
	p.popTillStop()
	out, ok := p.fragmentsStack.Pop()
	if !ok {
		return nil, errors.New("no fragments to return")
	}
	return &out, combineErrors(p.Errors)
}
