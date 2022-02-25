package compiler

import (
	"testing"

	"github.com/wingerjc/tableman-golang/pkg/program"
)

func TestBasicRolls(t *testing.T) {
	p, assert := setupParser(t)
	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext().
		SetRandom(rand)

	rand.AddMore(1)
	expr := `{ 1d1.mode? }`
	res := shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(1, res, assert)

	rand.AddMore(3)
	expr = `{ 1d5? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(3, res, assert)

	rand.AddMore(1, 2, 3)
	expr = `{ 3d9? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(6, res, assert)

	rand.AddMore(1, 2, 3, 6, 4)
	expr = `{ 5d10.sum? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(16, res, assert)
}

func TestRollFnAggs(t *testing.T) {
	p, assert := setupParser(t)
	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext().
		SetRandom(rand)

	rand.AddMore(4, 3, 2, 2, 3, 3)
	expr := `{ 6d8.mode? }`
	res := shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(3, res, assert)

	rand.AddMore(2, 1, 1)
	expr = `{ 3d2.mode? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(1, res, assert)

	rand.AddMore(2, 10, 9)
	expr = `{ 3d12.max? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(10, res, assert)

	rand.AddMore(17, 3, 2)
	expr = `{ 3d20.min? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(2, res, assert)

	rand.AddMore(4, 6, 2, 8)
	expr = `{ 4d8.avg? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(5, res, assert)

	rand.AddMore(8, 11)
	expr = `{ 2d20.avg? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(9, res, assert)

	rand.AddMore(11, 7, 19)
	expr = `{ 3d20.median? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(11, res, assert)

	rand.AddMore(1, 7, 10, 12)
	expr = `{ 4d20.median? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(8, res, assert)
}

func TestRollCountAggrs(t *testing.T) {
	p, assert := setupParser(t)
	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext().
		SetRandom(rand)

	rand.AddMore(4, 3, 2, 2, 3, 3)
	expr := `{ 6d8.+2.+3? }`
	res := shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(5, res, assert)

	rand.AddMore(1, 7, 10, 12)
	expr = `{ 4d20.-7.-10? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(-2, res, assert)

	rand.AddMore(3, 4, 6, 18)
	expr = `{ 4d20.-4x3.+18x6.+1x100? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(3, res, assert)
}

func TestRollDropValues(t *testing.T) {
	p, assert := setupParser(t)
	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext().
		SetRandom(rand)

	rand.AddMore(4, 3, 2, 2, 3, 3)
	expr := `{ 6d8l4? }`
	res := shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(10, res, assert)

	rand.AddMore(4, 3, 2, 2, 3, 3)
	expr = `{ 6d8h1? }`
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertInt(4, res, assert)
}

func TestRollStringValue(t *testing.T) {
	p, assert := setupParser(t)
	rand := program.NewTestRandSource()
	ctx := program.NewRootExecutionContext().
		SetRandom(rand)

	rand.AddMore(4, 3, 2, 2, 3, 3)
	expr := `{ 6d8l4.str? }`
	expect := "10: 6d8 (2, 2, 3, 3) drop(3, 4)"
	res := shouldParseExprWithContext(expr, p, ctx, assert)
	assertString(expect, res, assert)
	assert.Equal(expect, ctx.LatestRoll())

	rand.AddMore(4, 6, 1)
	expr = `{ 3d12.str? }`
	expect = "11: 3d12 (1, 4, 6)"
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertString(expect, res, assert)
	assert.Equal(expect, ctx.LatestRoll())

	rand.AddMore(3, 8, 4)
	expr = `{ 3d12.avg.str? }`
	expect = "5: 3d12 avg(3, 4, 8)"
	res = shouldParseExprWithContext(expr, p, ctx, assert)
	assertString(expect, res, assert)
	assert.Equal(expect, ctx.LatestRoll())
}
