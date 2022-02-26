package program

import (
	"fmt"
)

const (
	RootPack = "_ROOT"
)

type Evallable interface {
	Eval() ExpressionEval
}

type ExpressionEval interface {
	SetContext(ctx *ExecutionContext) ExpressionEval
	HasNext() bool
	Next() (ExpressionEval, error)
	Provide(res *ExpressionResult) error
	Resolve() (*ExpressionResult, error)
}

type Program struct {
	packs TableMap
	ctx   *ExecutionContext
}

func NewProgram(packs TableMap) *Program {
	ctx := NewRootExecutionContext()
	ctx.packs = packs
	return &Program{
		packs: packs,
		ctx:   ctx,
	}
}

func (p *Program) PackCount() int {
	return len(p.packs)
}

func (p *Program) Eval(expr Evallable) (*ExpressionResult, error) {
	return EvaluateExpression(expr, p.ctx.Child())
}

func NewTablePack(key string, name string, tables map[string]*Table) *TablePack {
	return &TablePack{
		key:    key,
		name:   name,
		tables: tables,
	}
}

type TablePack struct {
	key    string
	name   string
	tables map[string]*Table
}

type ResultType int

type ExpressionResult struct {
	resultType ResultType
	strVal     string
	intVal     int
}

func (e *ExpressionResult) Equal(other *ExpressionResult) bool {
	if e.resultType != other.resultType {
		return false
	}
	if e.resultType == IntResult {
		return e.intVal == other.intVal
	}
	if e.resultType == StringResult {
		return e.strVal == other.strVal
	}
	return false
}

func (e *ExpressionResult) BoolVal() bool {
	return e.resultType == IntResult && e.intVal != 0
}

func (e *ExpressionResult) SameType(other *ExpressionResult) bool {
	return e.resultType == other.resultType
}

func NewStringResult(val string) *ExpressionResult {
	return &ExpressionResult{
		resultType: StringResult,
		strVal:     val,
	}
}

func NewIntResult(val int) *ExpressionResult {
	return &ExpressionResult{
		resultType: IntResult,
		intVal:     val,
	}
}

func (e *ExpressionResult) MatchType(types ...ResultType) bool {
	for _, t := range types {
		if e.resultType == t {
			return true
		}
	}
	return false
}

func (e *ExpressionResult) IntVal() int {
	return e.intVal
}

func (e *ExpressionResult) StringVal() string {
	return e.strVal
}

const (
	AnyTypeResult = 0
	StringResult  = 1
	IntResult     = 2
)

type RollHistory struct {
	rollResults []string
}

func (h *RollHistory) ClearRolls() {
	h.rollResults = make([]string, 0)
}

func (h *RollHistory) GetRollHistory() []string {
	return h.rollResults[:]
}

func (h *RollHistory) AddRollToHistory(roll string) {
	h.rollResults = append(h.rollResults, roll)
}

func (h *RollHistory) LatestRoll() string {
	if len(h.rollResults) == 0 {
		return ""
	}
	return h.rollResults[len(h.rollResults)-1]
}

type ExecutionContext struct {
	*RollHistory
	parent *ExecutionContext
	values map[string]*ExpressionResult
	packs  TableMap
	rand   RandomSource
}

func NewRootExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		RollHistory: &RollHistory{
			rollResults: make([]string, 0),
		},
		parent: nil,
		values: make(map[string]*ExpressionResult),
		rand:   &DefaultRandSource{},
	}
}

func (ctx *ExecutionContext) SetRandom(r RandomSource) *ExecutionContext {
	ctx.rand = r
	return ctx
}

func (ctx *ExecutionContext) Rand(low int, high int) int {
	return ctx.rand.Get(low, high)
}

func (ctx *ExecutionContext) Child() *ExecutionContext {
	return &ExecutionContext{
		RollHistory: ctx.RollHistory,
		parent:      ctx,
		values:      make(map[string]*ExpressionResult),
		packs:       ctx.packs,
		rand:        ctx.rand,
	}
}

func (ctx *ExecutionContext) SetPacks(packs TableMap) {
	ctx.packs = packs
}

func (ctx *ExecutionContext) Set(key string, val *ExpressionResult) {
	ctx.values[key] = val
}

func (ctx *ExecutionContext) Resolve(key string) *ExpressionResult {
	for c := ctx; c != nil; c = c.parent {
		if v, ok := c.values[key]; ok {
			return v
		}
	}
	return nil
}

func EvaluateExpression(e Evallable, ctx *ExecutionContext) (*ExpressionResult, error) {
	if ctx == nil {
		ctx = NewRootExecutionContext()
	}
	stack := make([]ExpressionEval, 0)
	stack = append(stack, e.Eval().SetContext(ctx.Child()))
	for len(stack) > 0 {
		// See if we need to push another resolution node on the current stack.
		cur := stack[len(stack)-1]
		if cur.HasNext() {
			next, err := cur.Next()
			if err != nil {
				return nil, err
			}
			stack = append(stack, next)
			continue
		}
		result, err := cur.Resolve()
		if err != nil {
			return nil, err
		}
		if len(stack) == 1 {
			return result, nil
		}
		stack = stack[:len(stack)-1]
		if err = stack[len(stack)-1].Provide(result); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("ASDFASDFASDFASDF")
}

type TableMap map[string]*TablePack

func NewTableMap() TableMap {
	return make(map[string]*TablePack)
}
