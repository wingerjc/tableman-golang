package program

import (
	"fmt"
	"sort"
	"strings"
)

type Roll struct {
	print      bool
	diceCount  int
	diceSides  int
	selector   *RollSelect
	aggrFn     string
	countAggrs []*RollCountAggr
	rand       RandomSource
}

func NewRoll(count int, sides int, rand RandomSource) *Roll {
	return &Roll{
		diceCount:  count,
		diceSides:  sides,
		selector:   nil,
		aggrFn:     "",
		countAggrs: make([]*RollCountAggr, 0),
		print:      false,
		rand:       rand,
	}
}

func (r *Roll) Set() Evallable {
	return r
}

func (r *Roll) WithAggr(aggrFn string) *Roll {
	r.aggrFn = aggrFn
	return r
}

func (r *Roll) WithCountAggr(countAggrs []*RollCountAggr) *Roll {
	r.countAggrs = countAggrs
	return r
}

func (r *Roll) WithSelector(selector *RollSelect) *Roll {
	r.selector = selector
	return r
}

func (r *Roll) WithPrint(print bool) *Roll {
	r.print = print
	return r
}

func (r *Roll) Eval() ExpressionEval {
	return &rollEval{
		def: r,
	}
}

type RollSelect struct {
	high  bool
	count int
}

func NewRollSelect(isHigh bool, count int) *RollSelect {
	return &RollSelect{
		high:  isHigh,
		count: count,
	}
}

type RollCountAggr struct {
	number     int
	multiplier int
}

func NewRollCountAggr(number int, multiplier int) *RollCountAggr {
	return &RollCountAggr{
		number:     number,
		multiplier: multiplier,
	}
}

type rollEval struct {
	ctx *ExecutionContext
	def *Roll
}

func (r *rollEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	r.ctx = ctx
	return r
}

func (r *rollEval) HasNext() bool {
	return false
}

func (r *rollEval) Next() (ExpressionEval, error) {
	return nil, fmt.Errorf("roll has no sub-expressions")
}

func (r *rollEval) Provide(res *ExpressionResult) error {
	return fmt.Errorf("no values should be provided to roll expressions")
}

func (r *rollEval) Resolve() (*ExpressionResult, error) {
	res := &rollResult{
		value: 0,
		keep:  make([]int, 0),
		drop:  make([]int, 0),
	}
	for i := 0; i < r.def.diceCount; i++ {
		res.keep = append(res.keep, r.def.rand.Get(1, r.def.diceSides+1))
	}
	sort.Ints(res.keep)
	if r.def.selector != nil {
		toDrop := r.def.diceCount - r.def.selector.count
		if toDrop < 0 {
			return nil, fmt.Errorf("cannot drop more dice than rolled dice %d dropped %d",
				r.def.diceCount,
				r.def.selector.count,
			)
		}
		if r.def.selector.high {
			res.drop = res.keep[:toDrop]
			res.keep = res.keep[toDrop:]
		} else {
			res.drop = res.keep[r.def.selector.count:]
			res.keep = res.keep[:r.def.selector.count]
		}
	}

	// Can't have count aggregation and anything but string aggregator.
	if len(r.def.countAggrs) > 0 && len(r.def.aggrFn) > 0 && r.def.aggrFn != "roll" {
		return nil, fmt.Errorf("count aggregation can't be used with %s aggregation", r.def.aggrFn)
	}

	// Calculate avlue for count aggregations.
	for _, v := range res.keep {
		for _, aggr := range r.def.countAggrs {
			if v == aggr.number {
				res.value += aggr.multiplier
			}
		}
	}

	switch r.def.aggrFn {
	case "":
		if len(r.def.countAggrs) > 0 {
			break
		}
		fallthrough
	case "sum":
		for _, v := range res.keep {
			res.value += v
		}
	case "mode":
		counts := make(map[int]int)
		var cur int
		var ok bool
		for _, v := range res.keep {
			if cur, ok = counts[v]; ok {
				counts[v] = cur + 1
			} else {
				counts[v] = 1
			}
		}
		max := -1
		for k, v := range counts {
			if v > max {
				max = v
				cur = k
			}
		}
		res.value = cur
	case "avg":
		sum := 0
		for _, v := range res.keep {
			sum += v
		}
		res.value = sum / r.def.diceCount
	case "min":
		min := r.def.diceSides + 1
		for _, v := range res.keep {
			if v < min {
				min = v
			}
		}
		res.value = min
	case "max":
		max := -1
		for _, v := range res.keep {
			if v > max {
				max = v
			}
		}
		res.value = max
	default:
		return nil, fmt.Errorf("no roll aggregator matches '%s'", r.def.aggrFn)
	}
	if r.def.print {
		return printResult(r.def, res)
	}
	return NewIntResult(res.value), nil
}

func printResult(def *Roll, res *rollResult) (*ExpressionResult, error) {
	keepStr := make([]string, 0)
	for _, v := range res.keep {
		c := ""
		first := false
		for _, aggr := range def.countAggrs {
			if v == aggr.number {
				if first {
					c = "="
					first = false
				} else {
					c += ","
				}
				c += fmt.Sprintf("%d", aggr.multiplier)
			}
		}
		keepStr = append(keepStr, fmt.Sprintf("%d%s", v, c))
	}
	d := ""
	if len(res.drop) > 0 {
		dlist := make([]string, 0)
		for _, v := range res.drop {
			dlist = append(dlist, fmt.Sprintf("%d", v))
		}
		d = fmt.Sprintf("drop(%s)", strings.Join(dlist, ", "))
	}
	result := fmt.Sprintf(
		"%d: %dd%d %s(%s) %s",
		res.value,
		def.diceCount,
		def.diceSides,
		def.aggrFn,
		strings.Join(keepStr, ", "),
		d,
	)
	return NewStringResult(result), nil
}

type rollResult struct {
	value int
	keep  []int
	drop  []int
}
