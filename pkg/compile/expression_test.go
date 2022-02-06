package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
)

func TestExpression(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetExpressionParser()
	assert.NoError(err)

	expr := `{ foo }`
	result := shouldParseExpression(expr, p, assert)
	assertString("foo", result, assert)

	expr = `{ sum(6, 8) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(14, result, assert)

	// setting variables
	expr = `{ @foo=5; @foo }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(5, result, assert)

	// sequentially setting variables
	expr = `{ @foo=5, @bar=add(@foo,4), @baz=sub(@bar, 2); @baz }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(7, result, assert)
}
