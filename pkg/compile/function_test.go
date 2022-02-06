package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func assertInt(expect int, r *program.ExpressionResult, assert *assert.Assertions) {
	assert.True(r.MatchType(program.INT_RESULT))
	assert.Equal(expect, r.IntVal())
}

func assertString(expect string, r *program.ExpressionResult, assert *assert.Assertions) {
	assert.True(r.MatchType(program.STRING_RESULT))
	assert.Equal(expect, r.StringVal())
}

func shouldParseExpression(expr string, p *parser.ExpressionParser, assert *assert.Assertions) *program.ExpressionResult {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileExpression(parsed)
	assert.NoError(err)
	res, err := program.EvaluateExpression(prog)
	assert.NoError(err)
	return res
}

func TestAddFunc(t *testing.T) {
	assert := assert.New(t)
	// t.Parallel()
	p, _ := parser.GetExpressionParser()

	expr := `{ add(3, 7) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(10, result, assert)

	expr = `{ add(256) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(256, result, assert)

	expr = `{ add(3, 3, 3, 3, 3, 3, 3) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(21, result, assert)

	expr = `{ sum(1, 2, 3, 4, 0, 0, 0, 0, 0, 0, 0, 0) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(10, result, assert)

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

func TestSubFunc(t *testing.T) {
	assert := assert.New(t)
	// t.Parallel()
	p, err := parser.GetExpressionParser()
	assert.NoError(err)

	expr := `{ sub(10, 20) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(-10, result, assert)

	expr = `{ sub(21, 3, 3, 3, 3, 3, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ sub(6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(6, result, assert)

	expr = `{ sub(6, 8) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(-2, result, assert)

	expr = `{ sub(21, 3, 3, 3, 3, 3, 6, 0, 0, 0, 0, 0) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

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

func TestConcatFunc(t *testing.T) {
	assert := assert.New(t)
	// t.Parallel()
	p, err := parser.GetExpressionParser()
	assert.NoError(err)

	expr := `{ concat( foo ) }`
	result := shouldParseExpression(expr, p, assert)
	assertString("foo", result, assert)

	expr = `{ concat( foo, bar, "baz" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertString("foobarbaz", result, assert)

	expr = `{ @sp=" "; concat(that, @sp, sounds, @sp, "right!")}`
	result = shouldParseExpression(expr, p, assert)
	assertString("that sounds right!", result, assert)

	// At least 1 parameter, compile error
	expr = `{ concat() }`
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	_, err = CompileExpression(parsed)
	assert.Error(err)

	// Need all string values, runtime error
	expr = `{ concat( 5 ) }`
	parsed, err = p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileExpression(parsed)
	assert.NoError(err)
	_, err = program.EvaluateExpression(prog)
	assert.Error(err)

}
