package visitor

import (
	"errors"
	"fmt"

	"github.com/lejabque/software-design-2022/parser/internal/tokens"
)

type CalcVisitor struct {
	stack []int64
}

var ErrNotEnoughOperands = errors.New("not enough operands in expression")
var ErrIncorrectExpression = errors.New("incorrect expression")

func (c *CalcVisitor) VisitNumber(token tokens.Token) error {
	c.stack = append(c.stack, token.(tokens.Number).Value)
	return nil
}

func (c *CalcVisitor) VisitParen(token tokens.Token) error {
	return errors.New("Not supported in CalcVisitor, RPN only")
}

func (c *CalcVisitor) VisitOperation(token tokens.Token) error {
	if len(c.stack) < 2 {
		return ErrNotEnoughOperands
	}

	a := c.stack[len(c.stack)-2]
	b := c.stack[len(c.stack)-1]

	c.stack = c.stack[:len(c.stack)-2]

	switch token.(type) {
	case tokens.Plus:
		c.stack = append(c.stack, a+b)
	case tokens.Minus:
		c.stack = append(c.stack, a-b)
	case tokens.Multiply:
		c.stack = append(c.stack, a*b)
	case tokens.Divide:
		c.stack = append(c.stack, a/b)
	default: // should never happen
		panic(fmt.Errorf("unknown operation type %T", token))
	}

	return nil
}

func (c *CalcVisitor) VisitMultiple(tokens []tokens.Token) (int64, error) {
	for _, token := range tokens {
		if err := token.Accept(c); err != nil {
			return 0, err
		}
	}

	if len(c.stack) != 1 {
		return 0, ErrIncorrectExpression
	}

	return c.stack[0], nil
}
