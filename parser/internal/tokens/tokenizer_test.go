package tokens

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	var tokenizer Tokenizer
	tokens, err := tokenizer.Tokenize("(23 + 10) * 5 - 3 * (32 + 5) * (10 - 4 * 5) + 8 / 2")
	assert.NoError(t, err)
	expected := []Token{
		LeftParen{},
		Number{Value: 23},
		Plus{},
		Number{Value: 10},
		RightParen{},
		Multiply{},
		Number{Value: 5},
		Minus{},
		Number{Value: 3},
		Multiply{},
		LeftParen{},
		Number{Value: 32},
		Plus{},
		Number{Value: 5},
		RightParen{},
		Multiply{},
		LeftParen{},
		Number{Value: 10},
		Minus{},
		Number{Value: 4},
		Multiply{},
		Number{Value: 5},
		RightParen{},
		Plus{},
		Number{Value: 8},
		Divide{},
		Number{Value: 2},
	}
	assert.Equal(t, expected, tokens)
}
