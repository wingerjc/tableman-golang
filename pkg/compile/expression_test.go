package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func shouldParseExpression(expr string, p *parser.ExpressionParser, assert *assert.Assertions) *program.ExpressionResult {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileExpression(parsed)
	assert.NoError(err)
	res, err := program.EvaluateExpression(prog)
	assert.NoError(err)
	return res
}

func TestExpresion(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetExpressionParser()
	assert.NoError(err)

	expr := `{ foo }`
	res := shouldParseExpression(expr, p, assert)
	assert.Equal("foo", res.StringVal())

	expr = `{ sum(6, 8) }`
	res = shouldParseExpression(expr, p, assert)
	assert.Equal(14, res.IntVal())

	// setting variables
	expr = `{ @foo=5; @foo }`
	res = shouldParseExpression(expr, p, assert)
	assert.Equal(5, res.IntVal())

	// sequentially setting variables
	expr = `{ @foo=5, @bar=add(@foo,4), @baz=sub(@bar, 2); @baz }`
	res = shouldParseExpression(expr, p, assert)
	assert.Equal(7, res.IntVal())

	assert.Fail("foo")
}
