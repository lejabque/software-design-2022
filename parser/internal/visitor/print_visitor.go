package visitor

import (
	"io"

	"github.com/lejabque/software-design-2022/parser/internal/tokens"
)

type PrintVisitor struct {
	writer io.Writer
}

func NewPrintVisitor(writer io.Writer) *PrintVisitor {
	return &PrintVisitor{writer: writer}
}

func (p *PrintVisitor) Visit(token tokens.Token) error {
	_, err := p.writer.Write([]byte(token.String()))
	return err
}

func (p *PrintVisitor) VisitMultiple(tokens []tokens.Token) error {
	for i, token := range tokens {
		if err := token.Accept(p); err != nil {
			return err
		}
		if i != len(tokens)-1 {
			_, err := p.writer.Write([]byte(" "))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
