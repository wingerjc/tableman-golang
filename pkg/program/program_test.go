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
