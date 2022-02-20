package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
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

func getTableAndExprParsers(assert *assert.Assertions) (*parser.TableParser, *parser.ExpressionParser) {
	t, err := parser.GetTableParser()
	assert.NoError(err)
	e, err := parser.GetExpressionParser()
	assert.NoError(err)
	return t, e
}

type tableCallTestContext struct {
	packKeys nameMap
	eCtx     *program.ExecutionContext
	eParser  *parser.ExpressionParser
	assert   *assert.Assertions
}

func assertTableCallStr(code string, expect string, ctx *tableCallTestContext) {
	ast, err := ctx.eParser.Parse(code)
	ctx.assert.NoError(err)
	expr, err := CompileExpression(ast, ctx.packKeys)
	ctx.assert.NoError(err)
	val, err := program.EvaluateExpression(expr, ctx.eCtx.Child())
	ctx.assert.NoError(err)
	ctx.assert.True(val.MatchType(program.STRING_RESULT))
	ctx.assert.Equal(expect, val.StringVal())
}

func assertTableCallRuntimeErr(code string, ctx *tableCallTestContext) {
	ast, err := ctx.eParser.Parse(code)
	ctx.assert.NoError(err)
	expr, err := CompileExpression(ast, ctx.packKeys)
	ctx.assert.NoError(err)
	_, err = program.EvaluateExpression(expr, ctx.eCtx.Child())
	ctx.assert.Error(err)
}

func assertTableCallCompileErr(code string, ctx *tableCallTestContext) {
	ast, err := ctx.eParser.Parse(code)
	ctx.assert.NoError(err)
	_, err = CompileExpression(ast, ctx.packKeys)
	ctx.assert.Error(err)
}

func TestTableCallCompile(t *testing.T) {
	assert := assert.New(t)
	tParser, eParser := getTableAndExprParsers(assert)

	packKeys := make(nameMap)
	packKeys["foo"] = "bar"

	tCode := `TableDef: color
	c=2 15 first: "red"`
	tParsed, err := tParser.Parse(tCode)
	assert.NoError(err)
	table, err := CompileTable(tParsed, packKeys)
	assert.NoError(err)

	tableMap := make(map[string]*program.Table)
	tableMap["color"] = table
	pack := program.NewTablePack("bar", "foo", tableMap)
	packMap := make(program.TableMap)
	packMap["bar"] = pack
	eCtx := program.NewRootExecutionContext()
	eCtx.SetPacks(packMap)
	eCtx.SetRandom(&program.DefaultRandSource{})

	ctx := &tableCallTestContext{
		packKeys: packKeys,
		eParser:  eParser,
		eCtx:     eCtx,
		assert:   assert,
	}

	code := `{ !foo.color() }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(roll) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(weighted) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(index, 15) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(label, first) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(deck) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(deck, shuffle) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(deck, no-shuffle) }`
	assertTableCallStr(code, "red", ctx)

	code = `{ !foo.color(deck) }`
	assertTableCallRuntimeErr(code, ctx)

	code = `{ !foo.color(nonesuch) }`
	assertTableCallRuntimeErr(code, ctx)

	code = `{ !foo.color(label) }`
	assertTableCallRuntimeErr(code, ctx)

	code = `{ !foo.color(label, a, b, c) }`
	assertTableCallCompileErr(code, ctx)

	code = `{ !foo.color(index) }`
	assertTableCallRuntimeErr(code, ctx)

	code = `{ !foo.color(index, 1, 2, 3) }`
	assertTableCallCompileErr(code, ctx)

	code = `{ !foo.color(deck, blergh) }`
	assertTableCallRuntimeErr(code, ctx)

	code = `{ !foo.color(deck, shuffle, nope) }`
	assertTableCallCompileErr(code, ctx)

	code = `{ !foo.color(4) }`
	assertTableCallRuntimeErr(code, ctx)

	code = `{ !foo.color(helicopter_helicopter) }`
	assertTableCallRuntimeErr(code, ctx)
}
