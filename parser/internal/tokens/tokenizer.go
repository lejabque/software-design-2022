package tokens

import (
	"errors"
	"fmt"
	"unicode"
)

const eof = rune(-1)

type Tokenizer struct {
	tokens []Token
	state  state
}

type state interface {
	handle(t *Tokenizer, c rune) error
}

type startState struct{}

func (s *startState) handle(t *Tokenizer, c rune) error {
	switch {
	case c == '(':
		t.tokens = append(t.tokens, LeftParen{})
	case c == ')':
		t.tokens = append(t.tokens, RightParen{})
	case c == '+':
		t.tokens = append(t.tokens, Plus{})
	case c == '-':
		t.tokens = append(t.tokens, Minus{})
	case c == '*':
		t.tokens = append(t.tokens, Multiply{})
	case c == '/':
		t.tokens = append(t.tokens, Divide{})
	case unicode.IsDigit(c):
		t.state = &numberState{}
		t.state.handle(t, c)
	case c == eof:
		t.state = &endState{}
	default:
		if !unicode.IsSpace(c) {
			return fmt.Errorf("unexpected character, got: %c", c)
		}
	}
	return nil
}

type numberState struct {
	number int64
}

func (n *numberState) handle(t *Tokenizer, c rune) error {
	if unicode.IsDigit(c) {
		n.number = n.number*10 + int64(c-'0')
		return nil
	}
	t.tokens = append(t.tokens, Number{n.number})
	t.state = &startState{}
	return t.state.handle(t, c)
}

type endState struct{}

func (e *endState) handle(t *Tokenizer, c rune) error {
	return errors.New("unexpected character, there should be no more characters")
}

func (t *Tokenizer) Tokenize(s string) ([]Token, error) {
	t.state = &startState{}
	t.tokens = t.tokens[:0]
	for _, c := range s {
		if err := t.state.handle(t, c); err != nil {
			return nil, err
		}
	}
	err := t.state.handle(t, eof)
	return t.tokens, err
}
