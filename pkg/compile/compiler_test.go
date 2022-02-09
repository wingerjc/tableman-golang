package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	p, err := c.CompileString(table)
	assert.NoError(err)
	assert.NotNil(p)
}
