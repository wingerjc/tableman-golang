package program

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/k0kubun/pp"
)

type RandomSource interface {
	Get(low int, high int) int
}

type DefaultRandSource struct {
}

func (r *DefaultRandSource) Get(low int, high int) int {
	return rand.Intn(high-low) + low
}

type TestingRandSource struct {
	vals []int
}

func (r *TestingRandSource) Get(low int, high int) int {
	result := r.vals[0]
	r.vals = r.vals[1:]
	return result
}

func (r *TestingRandSource) AddMore(vals ...int) {
	r.vals = append(r.vals, vals...)
}

func NewTestRandSource(val ...int) *TestingRandSource {
	return &TestingRandSource{
		vals: val,
	}
}

type Table struct {
	name         string
	tags         map[string]string
	rows         []*TableRow
	rowsByLabel  map[string]*TableRow
	rowsByRange  []*Range
	totalWeight  int
	totalCount   int
	currentCount int
	defaultRow   int
}

func NewTable(name string, tags map[string]string, rows []*TableRow) *Table {
	result := &Table{
		name:         name,
		tags:         tags,
		rows:         rows,
		rowsByLabel:  make(map[string]*TableRow),
		rowsByRange:  make([]*Range, 0),
		totalWeight:  0,
		totalCount:   0,
		currentCount: 0,
		defaultRow:   -1,
	}
	for i, r := range result.rows {
		if len(r.label) > 0 {
			result.rowsByLabel[r.Label()] = r
		}
		result.totalCount += r.Count()
		result.currentCount += r.Count()
		result.totalWeight += r.Weight()
		if r.isDefault {
			result.defaultRow = i
		}
		for _, rng := range r.Ranges() {
			rng.SetRow(r)
			result.rowsByRange = append(result.rowsByRange, rng)
		}
	}
	return result
}

func (t *Table) Roll() Evallable {
	return &rowFuture{
		fn: func(ctx *ExecutionContext) Evallable {
			index := ctx.Rand(0, len(t.rows))
			return t.rows[index].Value()
		},
	}
}

func (t *Table) WeightedRoll() Evallable {
	return &rowFuture{
		fn: func(ctx *ExecutionContext) Evallable {
			roll := ctx.Rand(0, t.totalWeight)
			i := 0
			for {
				cur := t.rows[i].Weight
				pp.Println(roll, cur)
				if t.rows[i].Weight() > roll {
					return t.rows[i].Value()
				}
				roll -= t.rows[i].Weight()
				i++
			}
		},
	}
}

func (t *Table) LabelRoll(key string) (Evallable, error) {
	r, ok := t.rowsByLabel[key]
	if !ok {
		if t.defaultRow >= 0 {
			return t.rows[t.defaultRow].Value(), nil
		}
		return nil, fmt.Errorf("in table '%s' no row labelled '%s' and no default row", t.name, key)
	}
	return r.Value(), nil
}

func (t *Table) DeckDraw() (Evallable, error) {
	if t.currentCount == 0 {
		return nil, fmt.Errorf("deck draw called too many times without shuffle on table '%s'", t.name)
	}

	return &rowFuture{
		fn: func(ctx *ExecutionContext) Evallable {
			roll := ctx.Rand(0, t.currentCount)
			for _, r := range t.rows {
				if r.currentCount == 0 {
					continue
				}
				if roll < r.currentCount {
					r.currentCount--
					t.currentCount--
					return r.Value()
				}
				roll -= r.currentCount
			}
			return nil
		},
	}, nil
}

func (t *Table) Shuffle() {
	for _, r := range t.rows {
		r.currentCount = r.count
	}
	t.currentCount = t.totalCount
}

func (t *Table) IndexRoll(key int) (Evallable, error) {
	for _, rng := range t.rowsByRange {
		if rng.InRange(key) {
			return rng.Row().Value(), nil
		}
	}
	if t.defaultRow >= 0 {
		return t.rows[t.defaultRow].Value(), nil
	}
	return nil, fmt.Errorf("in table '%s' no index %d and no default row set", t.name, key)
}

func (t *Table) Name() string {
	return t.name
}

func (t *Table) TotalCount() int {
	return t.totalCount
}

func (t *Table) TotalWeight() int {
	return t.totalWeight
}

func (t *Table) RowCount() int {
	return len(t.rows)
}

func (t *Table) Tag(tagName string) (string, bool) {
	v, ok := t.tags[tagName]
	return v, ok
}

func (t *Table) Default() (Evallable, error) {
	if t.defaultRow < 0 {
		return nil, fmt.Errorf("no default set for table '%s'", t.name)
	}
	return t.rows[t.defaultRow].Value(), nil
}

type rowFuture struct {
	fn func(ctx *ExecutionContext) Evallable
}

func (r *rowFuture) Eval() ExpressionEval {
	return r
}

func (r *rowFuture) SetContext(ctx *ExecutionContext) ExpressionEval {
	return r.fn(ctx).Eval().SetContext(ctx)
}

func (r *rowFuture) HasNext() bool {
	return false
}

func (r *rowFuture) Next() (ExpressionEval, error) {
	return nil, fmt.Errorf("table row future has no subexpressions")
}

func (r *rowFuture) Provide(res *ExpressionResult) error {
	return fmt.Errorf("can't provide value to table row future")
}

func (r *rowFuture) Resolve() (*ExpressionResult, error) {
	return nil, fmt.Errorf("can't resolve table row future")
}

type TableRow struct {
	label        string
	rangeVal     []*Range
	weight       int
	count        int
	currentCount int
	isDefault    bool
	value        Evallable
}

func NewTableRow(label string, rangeVal []*Range, weight int, count int, isDefault bool, value Evallable) *TableRow {
	return &TableRow{
		label:        label,
		rangeVal:     rangeVal,
		weight:       weight,
		count:        count,
		currentCount: count,
		isDefault:    isDefault,
		value:        value,
	}
}

func (r *TableRow) Default() bool {
	return r.isDefault
}

func (r *TableRow) Label() string {
	return r.label
}

func (r *TableRow) Count() int {
	return r.count
}

func (r *TableRow) Weight() int {
	return r.weight
}

func (r *TableRow) Ranges() []*Range {
	return r.rangeVal
}

func (r *TableRow) Value() Evallable {
	return r.value
}

type ListExpression struct {
	items []Evallable
}

func NewListExpression(items []Evallable) Evallable {
	return &ListExpression{
		items: items,
	}
}

func (l *ListExpression) Eval() ExpressionEval {
	return &listExpressionEval{
		config:  l,
		results: make([]*ExpressionResult, 0),
	}
}

type listExpressionEval struct {
	ctx     *ExecutionContext
	config  *ListExpression
	results []*ExpressionResult
	index   int
}

func (l *listExpressionEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	l.ctx = ctx
	return l
}

func (l *listExpressionEval) HasNext() bool {
	return l.index < len(l.config.items)
}

func (l *listExpressionEval) Next() (ExpressionEval, error) {
	if l.index >= len(l.config.items) {
		return nil, fmt.Errorf("engine trying to evaluate too many sub-expressions for row")
	}
	return l.config.items[l.index].Eval().SetContext(l.ctx.Child()), nil
}

func (l *listExpressionEval) Provide(res *ExpressionResult) error {
	if l.index >= len(l.config.items) {
		return fmt.Errorf("extra return value to table expression list")
	}
	l.results = append(l.results, res)
	l.index++
	return nil
}

func (l *listExpressionEval) Resolve() (*ExpressionResult, error) {
	result := ""
	for _, i := range l.results {
		if i.MatchType(STRING_RESULT) {
			result = result + i.StringVal()
			continue
		}
		result = result + strconv.Itoa(i.IntVal())
	}
	return NewStringResult(result), nil
}

type Range struct {
	low  int
	high int
	row  *TableRow
}

func NewRange(low int, high int) *Range {
	return &Range{
		low:  low,
		high: high,
	}
}

func (r *Range) SetRow(row *TableRow) {
	r.row = row
}

func (r *Range) Row() *TableRow {
	return r.row
}

func (r *Range) InRange(val int) bool {
	return val <= r.high && val >= r.low
}
