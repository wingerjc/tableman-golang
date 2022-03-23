package program

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicRoll(t *testing.T) {
	assert := assert.New(t)
	random := NewTestRandSource()
	ctx := NewRootExecutionContext().
		SetRandom(random)

	random.AddMore(3, 4)
	r, err := NewRoll(2, 5).Eval().SetContext(ctx).Resolve()
	assert.NoError(err)
	assert.True(r.MatchType(IntResult))
	assert.Equal(7, r.IntVal())

	random.AddMore(5, 6, 3, 2, 5, 1)
	r, err = NewRoll(6, 10).
		WithAggr("mode").
		WithPrint(true).
		WithSelector(NewRollSelect(true, 4)).Eval().SetContext(ctx).Resolve()
	assert.NoError(err)
	assert.True(r.MatchType(StringResult))
	assert.Equal("5: 6d10 mode(3, 5, 5, 6) drop(1, 2)", r.StringVal())
}
