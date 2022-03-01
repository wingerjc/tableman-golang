package program

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/k0kubun/pp"
)

// Table is a program unit that can randomly and deterministically return
// row Evallable objects. The core of tableman.
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
	deckMu       sync.Mutex
}

// Copy deep copies a Table
func (t *Table) Copy() *Table {
	newRows := make([]*TableRow, 0, len(t.rows))
	for _, r := range t.rows {
		newRows = append(newRows, r.Copy())
	}
	return NewTable(t.name, t.tags, newRows)
}

// NewTable creates a new table object.
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
			rng.setRow(r)
			result.rowsByRange = append(result.rowsByRange, rng)
		}
	}
	return result
}

// Roll randomly on the table treating each row with equal weight.
func (t *Table) Roll() Evallable {
	return &rowFuture{
		fn: func(ctx *ExecutionContext) Evallable {
			index := ctx.Rand(0, len(t.rows))
			return t.rows[index].Value()
		},
	}
}

// WeightedRoll randomly rolls on the table using the defined row weights to
// decide which rows to return. Rows without set weights are treated as w=1.
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

// LabelRoll fetches a row directly from the table using the passed label.
// If no label was defined on the table, the default row will be returned.
// If no default row was specified, an error will be returned.
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

// DeckDraw will treat the table as a deck of cards using the count value (default 1) to
// choose a row from the deck.
//
// This is the only stateful draw from the table, once all counts have been exhausted
// an error will be returned if DeckDraw is called. Reset the counts by calling `Shuffle()`.
func (t *Table) DeckDraw() (Evallable, error) {
	t.deckMu.Lock()
	defer t.deckMu.Unlock()
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

// Shuffle resets all counts for DeckDraw calls.
func (t *Table) Shuffle() {
	for _, r := range t.rows {
		r.currentCount = r.count
	}
	t.currentCount = t.totalCount
}

// IndexRoll returns the row defined for the given index.
// If no index machtes, the default row will be used.
// If no default was defined, an error will be returned.
func (t *Table) IndexRoll(key int) (Evallable, error) {
	for _, rng := range t.rowsByRange {
		if rng.inRange(key) {
			return rng.getRow().Value(), nil
		}
	}
	if t.defaultRow >= 0 {
		return t.rows[t.defaultRow].Value(), nil
	}
	return nil, fmt.Errorf("in table '%s' no index %d and no default row set", t.name, key)
}

// Name returns the defined name of the table.
func (t *Table) Name() string {
	return t.name
}

// TotalCount returns the total number of "cards" for a DeckDraw.
func (t *Table) TotalCount() int {
	return t.totalCount
}

// TotalWeight returns the total weight of all tale rows.
func (t *Table) TotalWeight() int {
	return t.totalWeight
}

// RowCount returns the number of rows in the table.
func (t *Table) RowCount() int {
	return len(t.rows)
}

// Tag returns the tage valuve for the given tag name and a boolean for
// whether the tag was defined (the same way a map can).
func (t *Table) Tag(tagName string) (string, bool) {
	v, ok := t.tags[tagName]
	return v, ok
}

// Default returns the default row for the table or an error if ther isn't one.
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

// TableRow is an Evallable row for a tableman table.
type TableRow struct {
	label        string
	rangeVal     []*Range
	weight       int
	count        int
	currentCount int
	isDefault    bool
	value        Evallable
}

// NewTableRow creates a new TableRow object.
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

// Copy deep copies a TableRow
func (r *TableRow) Copy() *TableRow {
	return NewTableRow(
		r.label,
		r.rangeVal,
		r.weight,
		r.count,
		r.isDefault,
		r.value,
	)
}

// Default returns whether the row is a default value.
func (r *TableRow) Default() bool {
	return r.isDefault
}

// Label returns the label for the row, or an empty string if one isn't defined.
func (r *TableRow) Label() string {
	return r.label
}

// Count returns the current DeckDraw count for the row.
func (r *TableRow) Count() int {
	return r.count
}

// Weight returns the weight of the row.
func (r *TableRow) Weight() int {
	return r.weight
}

// Ranges returns the list of index ranges defined for this row.
func (r *TableRow) Ranges() []*Range {
	return r.rangeVal
}

// Value returns the Evalable avlue for this row.
func (r *TableRow) Value() Evallable {
	return r.value
}

// ListExpression is an Evallable that wraps all the items for a table row.
type ListExpression struct {
	items []Evallable
}

// NewListExpression creates a new list of epxressions for a row.
func NewListExpression(items []Evallable) Evallable {
	return &ListExpression{
		items: items,
	}
}

// Eval implementation for Evallable interface.
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
		if i.MatchType(StringResult) {
			result = result + i.StringVal()
			continue
		}
		result = result + strconv.Itoa(i.IntVal())
	}
	return NewStringResult(result), nil
}

// Range is a low-high number range.
//
// The range can be a single number if low and high are the same.
type Range struct {
	low  int
	high int
	row  *TableRow
}

// NewRange creates a new range value.
func NewRange(low int, high int) *Range {
	return &Range{
		low:  low,
		high: high,
	}
}

func (r *Range) setRow(row *TableRow) {
	r.row = row
}

func (r *Range) getRow() *TableRow {
	return r.row
}

func (r *Range) inRange(val int) bool {
	return val <= r.high && val >= r.low
}
