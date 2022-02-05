package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func TestAddfunc(t *testing.T) {
	assert := assert.New(t)
	// t.Parallel()
	p, _ := parser.GetExpressionParser()

	expr := `{ add(3, 7) }`
	result := shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(10, result.IntVal())

	expr = `{ add(256) }`
	result = shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(256, result.IntVal())

	expr = `{ add(3, 3, 3, 3, 3, 3, 3) }`
	result = shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(21, result.IntVal())

	expr = `{ sum(1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0) }`
	result = shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(10, result.IntVal())

	// Compile time error, not enough arguments
	expr = `{ add() }`
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	_, err = CompileValueExpr(parsed.Value)
	assert.Error(err)

	// runtime error, wrong argument type
	expr = `{ add( "foo" ) }`
	parsed, err = p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileValueExpr(parsed.Value)
	assert.NoError(err)
	_, err = program.EvaluateExpression(prog)
	assert.Error(err)
}

func TestSubfunc(t *testing.T) {
	assert := assert.New(t)
	// t.Parallel()
	p, err := parser.GetExpressionParser()
	assert.NoError(err)

	expr := `{ sub(10, 20) }`
	result := shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(-10, result.IntVal())

	expr = `{ sub(21, 3, 3, 3, 3, 3, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(0, result.IntVal())

	expr = `{ sub(6) }`
	result = shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(6, result.IntVal())

	expr = `{ sub(21, 3, 3, 3, 3, 3, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(0, result.IntVal())

	// Compile error, not enough arguments
	expr = `{ sub() }`
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	_, err = CompileValueExpr(parsed.Value)
	assert.Error(err)

	// Runtime error, wrong argument type(s)
	expr = `{ sub( bar, baz, qux ) }`
	parsed, err = p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileValueExpr(parsed.Value)
	assert.NoError(err)
	_, err = program.EvaluateExpression(prog)
	assert.Error(err)
}
