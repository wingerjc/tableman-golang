package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func parseRow(expr string, p *parser.RowParser, assert *assert.Assertions) *program.TableRow {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	row, err := CompileRow(parsed, DEFAULT_NAME_MAP)
	assert.NoError(err)
	return row
}

func parseTable(expr string, p *parser.TableParser, assert *assert.Assertions) *program.Table {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	table, err := CompileTable(parsed, DEFAULT_NAME_MAP)
	assert.NoError(err)
	return table
}

func TestTableRowCompile(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetRowParser()
	assert.NoError(err)

	expr := `"first"`
	row := parseRow(expr, p, assert)
	result, err := program.EvaluateExpression(row.Value(), nil)
	assert.NoError(err)
	assert.Equal("first", result.StringVal())
	assert.Equal(1, row.Count())
	assert.Equal(1, row.Weight())
	assert.Equal(false, row.Default())
	assert.Equal("", row.Label())
	assert.Len(row.Ranges(), 0)

	expr = `Default w=3 c=6 1,3-9,40 foo:{5}`
	row = parseRow(expr, p, assert)
	result, err = program.EvaluateExpression(row.Value(), nil)
	assert.NoError(err)
	assert.Equal("5", result.StringVal())
	assert.Equal(6, row.Count())
	assert.Equal(3, row.Weight())
	assert.Equal(true, row.Default())
	assert.Equal("foo", row.Label())
	assert.Len(row.Ranges(), 3)

	expr = `{foo} {bar} {baz}`
	row = parseRow(expr, p, assert)
	result, err = program.EvaluateExpression(row.Value(), nil)
	assert.NoError(err)
	assert.Equal("foobarbaz", result.StringVal())
}

func TestTableConstruction(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	expr := `TableDef: foo
	~ fruit: banana
	w=2 c=4: {2}
	Default w=3 c=9: {3}`
	table := parseTable(expr, p, assert)
	assert.Equal(2, table.RowCount())
	assert.Equal(5, table.TotalWeight())
	assert.Equal(13, table.TotalCount())
	tag, ok := table.Tag("fruit")
	assert.True(ok)
	assert.Equal("banana", tag)
	_, ok = table.Tag("hero")
	assert.False(ok)
	row, err := table.Default()
	assert.NoError(err)
	result, err := program.EvaluateExpression(row, nil)
	assert.NoError(err)
	assert.Equal("3", result.StringVal())

	expr = `TableDef: bar
	"first"
	"second"
	"third"`
	table = parseTable(expr, p, assert)
	assert.Equal(3, table.RowCount())
	assert.Equal(3, table.TotalWeight())
	assert.Equal(3, table.TotalCount())
	_, err = table.Default()
	assert.Error(err)
}

func TestTableRoll(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext()
	ctx.SetRandom(rand)

	expr := `TableDef: foo
	{1}
	1,6,7:{2}
	Default:{3}
	w=3:{4}
	c=19:{5}
	asdf:{6}`
	table := parseTable(expr, p, assert)
	rand.AddMore(5, 4, 3, 2, 1, 0)
	result, err := program.EvaluateExpression(table.Roll(), ctx)
	assert.NoError(err)
	assert.Equal("6", result.StringVal())
	result, err = program.EvaluateExpression(table.Roll(), ctx)
	assert.NoError(err)
	assert.Equal("5", result.StringVal())
	result, err = program.EvaluateExpression(table.Roll(), ctx)
	assert.NoError(err)
	assert.Equal("4", result.StringVal())
	result, err = program.EvaluateExpression(table.Roll(), ctx)
	assert.NoError(err)
	assert.Equal("3", result.StringVal())
	result, err = program.EvaluateExpression(table.Roll(), ctx)
	assert.NoError(err)
	assert.Equal("2", result.StringVal())
	result, err = program.EvaluateExpression(table.Roll(), ctx)
	assert.NoError(err)
	assert.Equal("1", result.StringVal())
}

func TestWeightedRoll(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext()
	ctx.SetRandom(rand)

	expr := `TableDef: foo
	w=3: {1}
	{2}
	Default w=4:{3}
	w=6 3-12:{4}
	c=19:{5}
	asdf:{6}`
	testCases := []int{1, 3, 14, 11}
	testExpect := []string{"1", "2", "5", "4"}
	table := parseTable(expr, p, assert)

	for i, val := range testCases {
		rand.AddMore(val)
		result, err := program.EvaluateExpression(table.WeightedRoll(), ctx)
		assert.NoError(err)
		assert.Equal(testExpect[i], result.StringVal())
	}
}

func TestLabelRoll(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext()
	ctx.SetRandom(rand)

	expr := `TableDef: foo
	w=3 once: {1}
	upon:{2}
	Default w=4 a:{3}
	w=6 3-12 time:{4}
	c=19 there:{5}
	"was a":{6}`
	testCases := []string{"was a", "time", "once", "there", "N/A"}
	testExpect := []string{"6", "4", "1", "5", "3"}
	table := parseTable(expr, p, assert)

	for i, val := range testCases {
		row, err := table.LabelRoll(val)
		assert.NoError(err)
		result, err := program.EvaluateExpression(row, ctx)
		assert.NoError(err)
		assert.Equal(testExpect[i], result.StringVal())
	}

	expr = `TableDef: bar
	{7}`
	table = parseTable(expr, p, assert)
	_, err = table.LabelRoll("anything  honestly")
	assert.Error(err)
}

func TestIndexRoll(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext()
	ctx.SetRandom(rand)

	expr := `TableDef: foo
	1,2,6-8: {1}
	Default: {2}
	w=4 13-15:{3}
	asdf:{4}
	9: {5}`
	testCases := []int{9, 8, 14, 128}
	testExpect := []string{"5", "1", "3", "2"}
	table := parseTable(expr, p, assert)

	for i, val := range testCases {
		row, err := table.IndexRoll(val)
		assert.NoError(err)
		result, err := program.EvaluateExpression(row, ctx)
		assert.NoError(err)
		assert.Equal(testExpect[i], result.StringVal())
	}

	expr = `TableDef: bar
	{999}`
	table = parseTable(expr, p, assert)
	_, err = table.IndexRoll(123)
	assert.Error(err)
}

func TestDeckRoll(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext()
	ctx.SetRandom(rand)

	expr := `TableDef: foo
	{1}
	Default w=3 c=2: {2}
	w=4:{3}
	asdf:{4}
	c=10 9: {5}`
	testCases := []int{0, 0, 12}
	testExpect := []string{"1", "2", "5"}
	table := parseTable(expr, p, assert)

	for i, val := range testCases {
		rand.AddMore(val)
		row, err := table.DeckDraw()
		assert.NoError(err)
		result, err := program.EvaluateExpression(row, ctx)
		assert.NoError(err)
		assert.Equal(testExpect[i], result.StringVal())
	}

	expr = `TableDef: bar
	{2}`
	table = parseTable(expr, p, assert)
	rand.AddMore(0, 0, 0, 0)
	row, err := table.DeckDraw()
	assert.NoError(err)
	_, err = program.EvaluateExpression(row, ctx)
	assert.NoError(err)
	_, err = table.DeckDraw()
	assert.Error(err)
	table.Shuffle()
	row, err = table.DeckDraw()
	assert.NoError(err)
	_, err = program.EvaluateExpression(row, ctx)
	assert.NoError(err)
}

func TestGeneratedTable(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetTableParser()
	assert.NoError(err)

	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext()
	ctx.SetRandom(rand)

	expr := `TableDef: foo
	["A","2","3","4", "5", "6", "7", "8", "9", "10" , "J" ,"Q",
	"K"][" of "]["Clubs", "Spades", "Diamonds", "Hearts"]`
	testCases := []int{0, 51, 14}
	testExpect := []string{"A of Clubs", "K of Hearts", "2 of Spades"}
	table := parseTable(expr, p, assert)
	assert.Equal(52, table.RowCount())
	assert.Equal(52, table.TotalCount())

	for i, val := range testCases {
		rand.AddMore(val)
		row := table.Roll()
		assert.NoError(err)
		result, err := program.EvaluateExpression(row, ctx)
		assert.NoError(err)
		assert.Equal(testExpect[i], result.StringVal())
	}

	row, err := table.IndexRoll(27)
	assert.NoError(err)
	result, err := program.EvaluateExpression(row, ctx)
	assert.NoError(err)
	assert.Equal("A of Diamonds", result.StringVal())
}
