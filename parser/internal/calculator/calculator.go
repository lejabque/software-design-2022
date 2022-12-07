package calculator

import (
	"bytes"

	"github.com/lejabque/software-design-2022/parser/internal/tokens"
	"github.com/lejabque/software-design-2022/parser/internal/visitor"
)

type Calculator struct{}

func (c *Calculator) Calculate(input string) (string, int64, error) {
	tokenizer := tokens.Tokenizer{}
	calc := visitor.CalcVisitor{}
	parser := visitor.ParserVisitor{}
	buf := bytes.Buffer{}
	printer := visitor.NewPrintVisitor(&buf)

	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		return "", 0, err
	}
	rpn, err := parser.VisitMultiple(tokens)
	if err != nil {
		return "", 0, err
	}
	if err = printer.VisitMultiple(rpn); err != nil {
		return "", 0, err
	}
	res, err := calc.VisitMultiple(rpn)
	return buf.String(), res, err
}
