package program

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestBasicRoll(t *testing.T) {
	assert := assert.New(t)
	ctx := NewRootExecutionContext()
	random := NewTestRandSource()

	random.AddMore(3, 4)
	r, err := NewRoll(2, 5, random).Set().Eval().SetContext(ctx).Resolve()
	assert.NoError(err)
	assert.True(r.MatchType(INT_RESULT))
	assert.Equal(7, r.IntVal())

	random.AddMore(5, 6, 3, 2, 5, 1)
	r, err = NewRoll(6, 10, random).
		WithAggr("mode").
		WithPrint(true).
		WithSelector(NewRollSelect(true, 4)).Set().Eval().SetContext(ctx).Resolve()
	assert.NoError(err)
	assert.True(r.MatchType(STRING_RESULT))

	pp.Println(r.StringVal())

	assert.Fail("foo")
}
