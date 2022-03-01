package program

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutionContextInheritance(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := NewRootExecutionContext()
	assert.Len(root.values, 0)
	_, err := root.Resolve("anything")
	assert.Error(err)

	root.Set("foo", NewIntResult(5))
	res, err := root.Resolve("foo")
	assert.NoError(err)
	assert.True(res.MatchType(IntResult))
	assert.Equal(5, res.IntVal())

	c := root.Child().Child().Child()
	res, err = c.Resolve("foo")
	assert.NoError(err)
	assert.True(res.MatchType(IntResult))
	assert.Equal(5, res.IntVal())

	c.Set("foo", NewStringResult("bar"))
	res, err = c.Resolve("foo")
	assert.NoError(err)
	assert.True(res.MatchType(StringResult))
	assert.Equal("bar", res.strVal)
	rootRes, err := root.Resolve("foo")
	assert.NoError(err)
	assert.True(rootRes.MatchType(IntResult))
	assert.Equal(5, rootRes.IntVal())

	d := c.Child()
	res, err = d.Resolve("foo")
	assert.NoError(err)
	assert.True(res.MatchType(StringResult))
	assert.False(res.MatchType(IntResult))
	assert.Equal("bar", res.StringVal())
}

// TestCloneProgram verifies that cloning a program will create separate copy that
// does not share state with the original and wil not interfere.
func TestCloneProgram(t *testing.T) {
	assert := assert.New(t)

	r := NewTableRow("", make([]*Range, 0), 1, 5, false, NewString("abc", true))
	table := NewTable("test", make(map[string]string), []*TableRow{r})
	pack := NewTablePack("foo", "test-pack", map[string]*Table{"test": table})
	prog := NewProgram(TableMap{RootPack: pack})

	// Deck draw state. -----------------------------
	assert.Equal(1, prog.PackCount())
	assert.Equal(5, prog.packs[RootPack].tables["test"].currentCount)
	assert.Equal(5, prog.packs[RootPack].tables["test"].rows[0].currentCount)
	expr, err := NewTableCall(RootPack, "", "test", []Evallable{NewString("deck", true)})
	assert.NoError(err)
	// Make an initial call into the table to reduce the number of pulls below max.
	res, err := prog.Eval(expr)
	assert.NoError(err)
	assert.Equal("abc", res.StringVal())
	assert.Equal(4, prog.packs[RootPack].tables["test"].currentCount)
	assert.Equal(4, prog.packs[RootPack].tables["test"].rows[0].currentCount)

	// copy and verify new table freshness.
	p2 := prog.Copy()
	assert.Equal(1, p2.PackCount())
	assert.Equal(5, p2.packs[RootPack].tables["test"].currentCount)
	assert.Equal(5, p2.packs[RootPack].tables["test"].rows[0].currentCount)

	// Pull from the old table to make sure no forward references.
	res, err = prog.Eval(expr)
	assert.NoError(err)
	assert.Equal("abc", res.StringVal())
	assert.Equal(5, p2.packs[RootPack].tables["test"].currentCount)
	assert.Equal(5, p2.packs[RootPack].tables["test"].rows[0].currentCount)
	assert.Equal(3, prog.packs[RootPack].tables["test"].currentCount)
	assert.Equal(3, prog.packs[RootPack].tables["test"].rows[0].currentCount)

	// Pull from the new table and verify no back references.
	res, err = p2.Eval(expr)
	assert.NoError(err)
	assert.Equal("abc", res.StringVal())
	assert.Equal(4, p2.packs[RootPack].tables["test"].currentCount)
	assert.Equal(4, p2.packs[RootPack].tables["test"].rows[0].currentCount)
	assert.Equal(3, prog.packs[RootPack].tables["test"].currentCount)
	assert.Equal(3, prog.packs[RootPack].tables["test"].rows[0].currentCount)

	// Variable state ----------------------
	prog.ctx.Set("foo", NewStringResult("bar"))
	expr = NewVariable("foo")
	res, err = prog.Eval(expr)
	assert.NoError(err)
	assert.Equal("bar", res.StringVal())

	// Copy and verify original does not change.
	p3 := prog.Copy()
	res, err = prog.Eval(expr)
	assert.NoError(err)
	assert.Equal("bar", res.StringVal())

	// Copy should not have variable defined.
	_, err = p3.Eval(expr)
	assert.Error(err)

	// Verify that no edits get made back to the original
	// Kind of dumb, but maybe worth it? (when in doubt it's just a test case).
	p3.ctx.Set("foo", NewStringResult("Zum"))
	p3.ctx.Set("baz", NewStringResult("qux"))
	// New dest values
	res, err = p3.Eval(expr)
	assert.NoError(err)
	assert.Equal("Zum", res.StringVal())
	expr2 := NewVariable("baz")
	res, err = p3.Eval(expr2)
	assert.NoError(err)
	assert.Equal("qux", res.StringVal())
	// old pack values
	res, err = prog.Eval(expr)
	assert.NoError(err)
	assert.Equal("bar", res.StringVal())
	_, err = prog.Eval(expr2)
	assert.Error(err)
}
