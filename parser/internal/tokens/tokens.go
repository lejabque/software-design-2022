package tokens

import (
	"fmt"
)

type Token interface {
	String() string
	Accept(TokenVisitor) error
}

type TokenVisitor interface {
	// we don't have overloading in Go, so we need to have a separate method for each type
	VisitOperation(token Token) error
	VisitNumber(token Token) error
	VisitParen(token Token) error
}

type Plus struct{}

func (p Plus) String() string {
	return "+"
}

func (p Plus) Accept(v TokenVisitor) error {
	return v.VisitOperation(p)
}

type Minus struct{}

func (m Minus) String() string {
	return "-"
}

func (m Minus) Accept(v TokenVisitor) error {
	return v.VisitOperation(m)
}

type Multiply struct{}

func (m Multiply) String() string {
	return "*"
}

func (m Multiply) Accept(v TokenVisitor) error {
	return v.VisitOperation(m)
}

type Divide struct{}

func (d Divide) String() string {
	return "/"
}

func (d Divide) Accept(v TokenVisitor) error {
	return v.VisitOperation(d)
}

type Number struct {
	Value int64
}

func (n Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (n Number) Accept(v TokenVisitor) error {
	return v.VisitNumber(n)
}

type LeftParen struct{}

func (l LeftParen) String() string {
	return "("
}

func (l LeftParen) Accept(v TokenVisitor) error {
	return v.VisitParen(l)
}

type RightParen struct{}

func (r RightParen) String() string {
	return ")"
}

func (r RightParen) Accept(v TokenVisitor) error {
	return v.VisitParen(r)
}
