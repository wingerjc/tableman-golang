package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

var (
	DEFAULT_NAME_MAP = make(nameMap)
)

func assertInt(expect int, r *program.ExpressionResult, assert *assert.Assertions) {
	assert.True(r.MatchType(program.INT_RESULT))
	assert.Equal(expect, r.IntVal())
}

func assertString(expect string, r *program.ExpressionResult, assert *assert.Assertions) {
	assert.True(r.MatchType(program.STRING_RESULT))
	assert.Equal(expect, r.StringVal())
}

func assertCompFail(expr string, p *parser.ExpressionParser, assert *assert.Assertions) {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	_, err = CompileExpression(parsed, DEFAULT_NAME_MAP)
	assert.Error(err)
}

func assertRuntimeFail(expr string, p *parser.ExpressionParser, assert *assert.Assertions) {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileExpression(parsed, DEFAULT_NAME_MAP)
	assert.NoError(err)
	_, err = program.EvaluateExpression(prog, nil)
	assert.Error(err)
}

func shouldParseExpression(expr string, p *parser.ExpressionParser, assert *assert.Assertions) *program.ExpressionResult {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	prog, err := CompileExpression(parsed, DEFAULT_NAME_MAP)
	assert.NoError(err)
	res, err := program.EvaluateExpression(prog, nil)
	assert.NoError(err)
	return res
}

func setupParser(t *testing.T) (*parser.ExpressionParser, *assert.Assertions) {
	assert := assert.New(t)
	p, err := parser.GetExpressionParser()
	assert.NoError(err)
	return p, assert
}

func TestAddFunc(t *testing.T) {
	p, assert := setupParser(t)

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
	assertCompFail(expr, p, assert)

	// runtime error, wrong argument type
	expr = `{ add( "foo" ) }`
	assertRuntimeFail(expr, p, assert)
}

func TestSubFunc(t *testing.T) {
	p, assert := setupParser(t)

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
	assertCompFail(expr, p, assert)

	// Runtime error, wrong argument type(s)
	expr = `{ sub( bar, baz, qux ) }`
	assertRuntimeFail(expr, p, assert)
}

func TestConcatFunc(t *testing.T) {
	p, assert := setupParser(t)

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
	assertCompFail(expr, p, assert)

	// Need all string values, runtime error
	expr = `{ concat( 5 ) }`
	assertRuntimeFail(expr, p, assert)
}

func TestUpperFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ upper( foo ) }`
	result := shouldParseExpression(expr, p, assert)
	assertString("FOO", result, assert)

	expr = `{ upper("Hello World") }`
	result = shouldParseExpression(expr, p, assert)
	assertString("HELLO WORLD", result, assert)

	// compiler error, min 1 param
	expr = `{ upper() }`
	assertCompFail(expr, p, assert)

	// compiler error max 1 param
	expr = `{ upper(foo, bar) }`
	assertCompFail(expr, p, assert)

	// runtime error, string only
	expr = `{ upper( 7 ) }`
	assertRuntimeFail(expr, p, assert)
}

func TestLowerFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ lower( FOO ) }`
	result := shouldParseExpression(expr, p, assert)
	assertString("foo", result, assert)

	expr = `{ lower("hELLO wOrlD") }`
	result = shouldParseExpression(expr, p, assert)
	assertString("hello world", result, assert)

	// compiler error, min 1 param
	expr = `{ lower() }`
	assertCompFail(expr, p, assert)

	// compiler error max 1 param
	expr = `{ lower(foo, bar) }`
	assertCompFail(expr, p, assert)

	// runtime error, string only
	expr = `{ lower( 7 ) }`
	assertRuntimeFail(expr, p, assert)
}

func TestToStrFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ str( 6 ) }`
	result := shouldParseExpression(expr, p, assert)
	assertString("6", result, assert)

	expr = `{ str( -23 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertString("-23", result, assert)

	expr = `{ str( 5287 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertString("5287", result, assert)

	expr = `{ str( "123" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertString("123", result, assert)

	// compiler error, min 1 param
	expr = `{ str() }`
	assertCompFail(expr, p, assert)

	// compiler error max 1 param
	expr = `{ str( 123, 456) }`
	assertCompFail(expr, p, assert)
}

func TestToIntFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ int( "5" ) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(5, result, assert)

	expr = `{ int( "-8357" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(-8357, result, assert)

	expr = `{ int( 7 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(7, result, assert)

	// compiler error, min 1 param
	expr = `{ int() }`
	assertCompFail(expr, p, assert)

	// compiler error max 1 param
	expr = `{ int(foo, bar) }`
	assertCompFail(expr, p, assert)

	// runtime error, string parse error
	expr = `{ int( "not a number" ) }`
	assertRuntimeFail(expr, p, assert)
}

func TestEqFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ eq( 1, 2 ) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ eq( 101, 101 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ eq( "thing", "thing" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ eq( "thing", "other" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ eq( "thing", "other" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ eq( "thing", 8 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ eq( "thing", "other" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ eq() }`
	assertCompFail(expr, p, assert)

	expr = `{ eq(1, 1, 1) }`
	assertCompFail(expr, p, assert)
}

func TestGtFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ gt( 1, 2 ) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ gt( 1, 1 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ gt( 13, 1 ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ gt( "a", "a" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ gt( "a", "z" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ gt( "z", "a" ) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ gt() }`
	assertCompFail(expr, p, assert)

	expr = `{ gt( 1, 1, 1) }`
	assertCompFail(expr, p, assert)
}

func TestGteFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ gte(1, 2) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ gte(6, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ gte(14, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ gte("a", "z") }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ gte("a", "a") }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ gte(z, a) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ gte() }`
	assertCompFail(expr, p, assert)

	expr = `{ gte(1, 1, 1) }`
	assertCompFail(expr, p, assert)

	expr = `{ gte(a, 5) }`
	assertRuntimeFail(expr, p, assert)
}

func TestLtFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ lt(1, 2) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ lt(6, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ lt(14, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ lt("a", "z") }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ lt("a", "a") }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ lt(z, a) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ lt() }`
	assertCompFail(expr, p, assert)

	expr = `{ lt(1, 1, 1) }`
	assertCompFail(expr, p, assert)

	expr = `{ lt(a, 5) }`
	assertRuntimeFail(expr, p, assert)
}

func TestLteFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ lte(1, 2) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ lte(6, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ lte(14, 6) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ lte("a", "z") }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ lte("a", "a") }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ lte(z, a) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ lte() }`
	assertCompFail(expr, p, assert)

	expr = `{ lte(1, 1, 1) }`
	assertCompFail(expr, p, assert)

	expr = `{ lte(a, 5) }`
	assertRuntimeFail(expr, p, assert)
}

func TestAndFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ and(1, 1) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ and(1, 0) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ and(0, 0) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ and(1, 1, 1, 1, 1, 1) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ and(1, 1, 1, 1, 0, 1) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ and() }`
	assertCompFail(expr, p, assert)

	expr = `{ and( foo, 0) }`
	assertRuntimeFail(expr, p, assert)
}

func TestOrFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ or( 0, 0) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ or( 0, 1) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ or( 0, 1, 0, 1, 0, 1) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ or( 0, 0, 0, 0, 0) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ or() }`
	assertCompFail(expr, p, assert)

	expr = `{ or(1) }`
	assertCompFail(expr, p, assert)

	expr = `{ or(1, foo) }`
	assertRuntimeFail(expr, p, assert)
}

func TestNotFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ not(0) }`
	result := shouldParseExpression(expr, p, assert)
	assertInt(1, result, assert)

	expr = `{ not(4) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ not(-3) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ not(1) }`
	result = shouldParseExpression(expr, p, assert)
	assertInt(0, result, assert)

	expr = `{ not() }`
	assertCompFail(expr, p, assert)

	expr = `{ not( 1, 1) }`
	assertCompFail(expr, p, assert)

	expr = `{ not(foo) }`
	assertRuntimeFail(expr, p, assert)
}

func TestIfFunc(t *testing.T) {
	p, assert := setupParser(t)

	expr := `{ if(1, true, false) }`
	result := shouldParseExpression(expr, p, assert)
	assertString("true", result, assert)

	expr = `{ if(0, true, false) }`
	result = shouldParseExpression(expr, p, assert)
	assertString("false", result, assert)

	expr = `{ if() }`
	assertCompFail(expr, p, assert)

	expr = `{ if(1) }`
	assertCompFail(expr, p, assert)

	expr = `{ if(1, true) }`
	assertCompFail(expr, p, assert)

	expr = `{ if(1, true, false, other) }`
	assertCompFail(expr, p, assert)

	expr = `{ if(asdf, true, false) }`
	assertRuntimeFail(expr, p, assert)
}
