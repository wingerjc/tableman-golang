package compiler

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	4,5,6: { @foo=!first(); concat(@foo, !call(@foo) ) }
	`
	p, err := c.CompileString(table)
	assert.NoError(err)
	assert.NotNil(p)
}

func TestBasicFileRead(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := assert.New(t)
	dir := t.TempDir()
	pack1 := `TablePack: foo
	Import: f"%s" As: Other.Pack
	
	TableDef: first
	{1}
	-------------`
	pack2 := `TablePack: bar.baz
	Import: f"%s" As: Last.pack
	
	TableDef: second
	{2}`
	pack3 := `TablePack: qux
	Import: f"%s"
	
	TableDef: third
	{3}`

	// Double import with loop
	f3Name := filepath.Join(dir, "f3")
	f2Name := filepath.Join(dir, "f2")
	err := ioutil.WriteFile(f3Name, []byte(fmt.Sprintf(pack3, f2Name)), 0644)
	assert.NoError(err)
	err = ioutil.WriteFile(f2Name, []byte(fmt.Sprintf(pack2, f3Name)), 0644)
	assert.NoError(err)

	code := fmt.Sprintf(pack1, f2Name)
	c, err := NewCompiler()
	assert.NoError(err)
	prog, err := c.CompileString(code)
	assert.NoError(err)
	assert.Equal(4, prog.PackCount())

	// Reading from
	f1Name := filepath.Join(dir, "f1")
	err = ioutil.WriteFile(f1Name, []byte(code), 0644)
	assert.NoError(err)
	prog, err = c.CompileFile(f1Name)
	assert.NoError(err)
	assert.Equal(4, prog.PackCount())
}

func TestImportedTableCalls(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	assert := assert.New(t)
	dir := t.TempDir()

	assert.NotZero(len(dir))
}
