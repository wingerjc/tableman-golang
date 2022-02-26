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
	assert.Nil(root.Resolve("anything"))

	root.Set("foo", NewIntResult(5))
	res := root.Resolve("foo")
	assert.NotNil(res)
	assert.True(res.MatchType(IntResult))
	assert.Equal(5, res.IntVal())

	c := root.Child().Child().Child()
	res = c.Resolve("foo")
	assert.NotNil(res)
	assert.True(res.MatchType(IntResult))
	assert.Equal(5, res.IntVal())

	c.Set("foo", NewStringResult("bar"))
	res = c.Resolve("foo")
	assert.NotNil("foo")
	assert.True(res.MatchType(StringResult))
	assert.Equal("bar", res.strVal)
	rootRes := root.Resolve("foo")
	assert.True(rootRes.MatchType(IntResult))
	assert.Equal(5, rootRes.IntVal())

	d := c.Child()
	res = d.Resolve("foo")
	assert.NotNil(res)
	assert.True(res.MatchType(StringResult))
	assert.False(res.MatchType(IntResult))
	assert.Equal("bar", res.StringVal())
}
