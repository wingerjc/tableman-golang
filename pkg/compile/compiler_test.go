package compiler

import (
	"testing"

	"github.com/k0kubun/pp"
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
	p, errs := c.CompileString(table)
	for _, e := range errs {
		pp.Println(e)
	}
	assert.Nil(errs)
	assert.NotNil(p)
}
