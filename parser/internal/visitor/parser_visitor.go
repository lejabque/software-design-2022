package visitor

import (
	"fmt"

	"github.com/lejabque/software-design-2022/parser/internal/tokens"
)

var priority = map[tokens.Token]int{
	tokens.Plus{}:     1,
	tokens.Minus{}:    1,
	tokens.Multiply{}: 2,
	tokens.Divide{}:   2,
}

// convert to RPN
type ParserVisitor struct {
	stack  []tokens.Token
	result []tokens.Token
}

func (p *ParserVisitor) VisitNumber(token tokens.Token) error {
	p.result = append(p.result, token)
	return nil
}

func (p *ParserVisitor) VisitParen(token tokens.Token) error {
	switch token.(type) {
	case tokens.LeftParen:
		p.stack = append(p.stack, token)
	case tokens.RightParen:
		for len(p.stack) > 0 {
			switch top := p.stack[len(p.stack)-1].(type) {
			case tokens.LeftParen:
				p.stack = p.stack[:len(p.stack)-1]
				return nil
			case tokens.RightParen:
				return fmt.Errorf("incorrect expression, operation or left parent expected, but got left parent: %v", token)
			case tokens.Number:
				return fmt.Errorf("incorrect expression, operation or left parent expected, but got a number: %v", token)
			default:
				p.result = append(p.result, top)
				p.stack = p.stack[:len(p.stack)-1]
			}
		}
	default:
		panic(fmt.Errorf("unknown brace type %T", token))
	}
	return nil
}

func (p *ParserVisitor) VisitOperation(token tokens.Token) error {
Loop:
	for len(p.stack) > 0 {
		switch top := p.stack[len(p.stack)-1].(type) {
		case tokens.LeftParen:
			break Loop
		case tokens.RightParen:
			break Loop
		case tokens.Number:
			break Loop
		default:
			if priority[top] >= priority[token] {
				p.result = append(p.result, top)
				p.stack = p.stack[:len(p.stack)-1]
			} else {
				break Loop
			}
		}
	}
	p.stack = append(p.stack, token)
	return nil
}

func (p *ParserVisitor) VisitMultiple(tokensToVisit []tokens.Token) ([]tokens.Token, error) {
	for _, token := range tokensToVisit {
		if err := token.Accept(p); err != nil {
			return nil, err
		}
	}
	for len(p.stack) > 0 {
		switch top := p.stack[len(p.stack)-1].(type) {
		case tokens.LeftParen:
			return nil, fmt.Errorf("incorrect expression, right parent expected, but got left parent: %v", top)
		case tokens.RightParen:
			return nil, fmt.Errorf("incorrect expression, operation or left parent expected, but got right parent: %v", top)
		case tokens.Number:
			return nil, fmt.Errorf("incorrect expression, operation or left parent expected, but got a number: %v", top)
		default:
			p.result = append(p.result, top)
			p.stack = p.stack[:len(p.stack)-1]
		}
	}
	return p.result, nil
}
