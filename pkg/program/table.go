package program

import (
	"fmt"
	"strconv"
)

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
	return nil
}

func (t *Table) WeightedRoll() Evallable {
	return nil
}

func (t *Table) DeckDraw() Evallable {
	return nil
}

func (t *Table) IndexRoll() Evallable {
	return nil
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

func (l *listExpressionEval) Next() ExpressionEval {
	return l.config.items[l.index].Eval().SetContext(l.ctx.Child())
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
