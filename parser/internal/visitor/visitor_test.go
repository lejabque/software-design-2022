package visitor

import (
	"bytes"
	"testing"

	"github.com/lejabque/software-design-2022/parser/internal/tokens"
	"github.com/stretchr/testify/assert"
)

func TestVisitors(t *testing.T) {
	var calc CalcVisitor
	var parser ParserVisitor

	input := []tokens.Token{
		tokens.LeftParen{},
		tokens.Number{Value: 23},
		tokens.Plus{},
		tokens.Number{Value: 10},
		tokens.RightParen{},
		tokens.Multiply{},
		tokens.Number{Value: 5},
		tokens.Minus{},
		tokens.Number{Value: 3},
		tokens.Multiply{},
		tokens.LeftParen{},
		tokens.Number{Value: 32},
		tokens.Plus{},
		tokens.Number{Value: 5},
		tokens.RightParen{},
		tokens.Multiply{},
		tokens.LeftParen{},
		tokens.Number{Value: 10},
		tokens.Minus{},
		tokens.Number{Value: 4},
		tokens.Multiply{},
		tokens.Number{Value: 5},
		tokens.RightParen{},
		tokens.Plus{},
		tokens.Number{Value: 8},
		tokens.Divide{},
		tokens.Number{Value: 2},
	}

	rpn, err := parser.VisitMultiple(input)
	assert.NoError(t, err)
	rpnExpected := []tokens.Token{
		tokens.Number{Value: 23},
		tokens.Number{Value: 10},
		tokens.Plus{},
		tokens.Number{Value: 5},
		tokens.Multiply{},
		tokens.Number{Value: 3},
		tokens.Number{Value: 32},
		tokens.Number{Value: 5},
		tokens.Plus{},
		tokens.Multiply{},
		tokens.Number{Value: 10},
		tokens.Number{Value: 4},
		tokens.Number{Value: 5},
		tokens.Multiply{},
		tokens.Minus{},
		tokens.Multiply{},
		tokens.Minus{},
		tokens.Number{Value: 8},
		tokens.Number{Value: 2},
		tokens.Divide{},
		tokens.Plus{},
	}

	assert.Equal(t, rpnExpected, rpn)

	// writer from string:
	buf := new(bytes.Buffer)
	printer := NewPrintVisitor(buf)
	err = printer.VisitMultiple(rpn)
	assert.NoError(t, err)
	assert.Equal(t, "23 10 + 5 * 3 32 5 + * 10 4 5 * - * - 8 2 / +", buf.String())

	res, err := calc.VisitMultiple(rpn)
	assert.NoError(t, err)
	assert.Equal(t, int64(1279), res)
}
