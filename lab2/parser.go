package main

import (
	"errors"
	"fmt"
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

func (p *Parser) BuildAST() (*Node, error) {
	tok, sym := p.lexer.Next()
	openParenths := 0

	// To insert '.' if wasn't
	// wasBinary := true
	for tok != EOS {
		switch tok {
		case Character:
			p.fragmentsStack.Push(&CharNode{Character: sym})
		case ClosedParenthesis:
			if openParenths < 0 {
				return nil, errors.New("unmatched ')' found")
			}
			p.popTillStop()
		}
		tok, sym = p.lexer.Next()
		p.fragmentsStack.Push(&CharNode{Character: 'a'})
	}
	res, _ := p.fragmentsStack.Pop()
	return &res, nil
}
