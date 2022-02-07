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
	row, err := CompileRow(parsed)
	assert.NoError(err)
	return row
}

func parseTable(expr string, p *parser.TableParser, assert *assert.Assertions) *program.Table {
	parsed, err := p.Parse(expr)
	assert.NoError(err)
	table, err := CompileTable(parsed)
	assert.NoError(err)
	return table
}

func TestTableRowCompile(t *testing.T) {
	assert := assert.New(t)
	p, err := parser.GetRowParser()
	assert.NoError(err)

	expr := `"first"`
	row := parseRow(expr, p, assert)
	result, err := program.EvaluateExpression(row.Value())
	assert.NoError(err)
	assert.Equal("first", result.StringVal())
	assert.Equal(1, row.Count())
	assert.Equal(1, row.Weight())
	assert.Equal(false, row.Default())
	assert.Equal("", row.Label())
	assert.Len(row.Ranges(), 0)

	expr = `Default w=3 c=6 1,3-9,40 foo:{5}`
	row = parseRow(expr, p, assert)
	result, err = program.EvaluateExpression(row.Value())
	assert.NoError(err)
	assert.Equal("5", result.StringVal())
	assert.Equal(6, row.Count())
	assert.Equal(3, row.Weight())
	assert.Equal(true, row.Default())
	assert.Equal("foo", row.Label())
	assert.Len(row.Ranges(), 3)

	expr = `{foo} {bar} {baz}`
	row = parseRow(expr, p, assert)
	result, err = program.EvaluateExpression(row.Value())
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
	result, err := program.EvaluateExpression(row)
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
