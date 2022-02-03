package compiler

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func TestCompileString(t *testing.T) {
	assert := assert.New(t)
	t.Parallel()

	c, err := NewCompiler()
	assert.NoError(err)

	table := `TablePack: foo
	
	TableDef: first
	"a"
	"b"
	"c"
	"d"
	
	TableDef: TheOtherOne
	Default 1,2,3: "mice"
	4,5,6: { @foo=!first(); concat(@foo, call(@foo) ) }
	`
	p, errs := c.CompileString(table)
	for _, e := range errs {
		pp.Println(e)
	}
	assert.Nil(errs)
	assert.NotNil(p)
}

func TestExpressionCompilation(t *testing.T) {
	assert := assert.New(t)
	// t.Parallel()
	p, _ := parser.GetExpressionParser()

	expr := `{ add(3, 7) }`
	parsed := &parser.Expression{}
	err := p.ParseString("", expr, parsed)
	assert.NoError(err)
	prog, err := CompileValueExpr(parsed.Value)
	assert.NoError(err)
	result, err := program.EvaluateExpression(prog)
	assert.NoError(err)
	assert.True(result.MatchType(program.INT_RESULT))
	assert.Equal(10, result.IntVal())
}
