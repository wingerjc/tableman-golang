package program

import (
	"fmt"
	"sync"
)

const (
	// RootPack is the default key for table calls.
	RootPack = "_ROOT"
)

// Evallable is an interface for a loaded program unit to provide
// an executable ExpressionEval during program execution.
//
// The returned ExpressionEval should always be semantically the same, but each value
// returned should have unique state (excepting deck draw counts).
//
// The root Evallable is should be considered a stateless generator for executable
// expressions.
//
//  Example:
//    An Evallable for the expression `if(eq(@foo, 3), "high", "low")` should always
//    return an ExpressionEval that acts the same for a given setting of @foo. But each
//    value returned may be provided a different value for @foo, and those different
//    settings should not interact or interfere with eachothers evaluation.
type Evallable interface {
	// Eval returns a stateful executable version of this expression.
	Eval() ExpressionEval
}

// ExpressionEval is an interface for evaluating any program node.
// It acts much lit an iterator over sub-expressions and allows for
// program nodes to use internal logic for sub-expression evaluation.
//
// This is assumed to be a stateful object.
//
// The main evaluation loop for a single ExpressionEval is similar to:
//  process(e) -> result:
//    while(e.HasNext):
//      e.Provide(process(e.Next))
//    return e.Resolve
// With extra care being taken around runtime errors generated.
//
// Implementations are allowed to skip providing sub-expressions or
// execute sub-expressions in any order,
// but it is recommended to use a right-to-left evaluation scheme as convention.
// For example: the `if(Cond, T, F)` flow control function will only evaluate
// T or F (depending on Cond's truthiness), but will always evaluate Cond first.
//
// Type are not enforced until evaluation so a given ExpressionEval can Resolve
// to any type necessary based on internal logic.
type ExpressionEval interface {
	// SetContext returns this ExpressionEval after
	SetContext(ctx *ExecutionContext) ExpressionEval

	// HasNext returns whether there is another sub-expression to evaluate.
	HasNext() bool

	// Next returns the next sub-expression to be evaluated.
	// Calling next multiple times should return the same semantic value to evaluate,
	// but is not required to be the same pointer/reference/memory value.
	Next() (ExpressionEval, error)

	// Provide sets the result for the sub-expresion currently returned by `Next()`.
	Provide(res *ExpressionResult) error

	// Resolve should be called to compute and return the final value of this
	// expression when `HasNext()` returns false.
	Resolve() (*ExpressionResult, error)
}

// Program is a set of TablePacks that can evaluate expressions as programs.
type Program struct {
	packs TableMap
	ctx   *ExecutionContext
}

// NewProgram creates a new program from a keyed set of tablepacks.
func NewProgram(packs TableMap) *Program {
	ctx := NewRootExecutionContext()
	ctx.packs = packs
	return &Program{
		packs: packs,
		ctx:   ctx,
	}
}

// PackCount returns the number of packs (files) loaded into this program.
func (p *Program) PackCount() int {
	return len(p.packs)
}

// Eval Evaluates a given Evallable against this program's state (tables+context).
func (p *Program) Eval(expr Evallable) (*ExpressionResult, error) {
	return EvaluateExpression(expr, p.ctx.Child())
}

// Copy returns a deep copy of the Program
func (p *Program) Copy() *Program {
	packs := make(TableMap)
	for k, v := range p.packs {
		packs[k] = v.Copy()
	}
	return NewProgram(packs)
}

// NewTablePack creates a new TablePack with the given tables.
func NewTablePack(key string, name string, tables map[string]*Table) *TablePack {
	return &TablePack{
		key:    key,
		name:   name,
		tables: tables,
	}
}

// TablePack represents a single executable tableman source file.
type TablePack struct {
	key    string
	name   string
	tables map[string]*Table
}

// Copy deep copies a TablePack
func (t *TablePack) Copy() *TablePack {
	tables := make(map[string]*Table)
	for k, v := range t.tables {
		tables[k] = v.Copy()
	}
	return NewTablePack(t.key, t.name, tables)
}

// ResultType is an alias for allowed return types from an expression.
type ResultType int

const (
	// AnyTypeResult can be used in type matching to take any result value.
	AnyTypeResult ResultType = 0

	// StringResult can be used to define or match to string expression results.
	StringResult ResultType = 1

	// IntResult can be used to define or matche to number/integer results.
	IntResult ResultType = 2
)

// ExpressionResult is a final result for an evaluated expression.
type ExpressionResult struct {
	resultType ResultType
	strVal     string
	intVal     int
}

// Equal compares 2 expression results for deep equality.
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

// BoolVal returns an integer-based boolean value for the result.
// 0 for false, any other integer for true.
// If the result is a string, false is always returned.
func (e *ExpressionResult) BoolVal() bool {
	return e.resultType == IntResult && e.intVal != 0
}

// SameType does a simple type comparison between two results.
func (e *ExpressionResult) SameType(other *ExpressionResult) bool {
	return e.resultType == other.resultType
}

// NewStringResult creates a new string-valued result.
func NewStringResult(val string) *ExpressionResult {
	return &ExpressionResult{
		resultType: StringResult,
		strVal:     val,
	}
}

// NewIntResult creates a new int-valued result.
func NewIntResult(val int) *ExpressionResult {
	return &ExpressionResult{
		resultType: IntResult,
		intVal:     val,
	}
}

// MatchType verifies that the type of this result is in the
// passed list. Setup for possible future types beyond int/string.
func (e *ExpressionResult) MatchType(types ...ResultType) bool {
	for _, t := range types {
		if e.resultType == t {
			return true
		}
	}
	return false
}

// IntVal returns the integer value, if it was set. Default 0.
func (e *ExpressionResult) IntVal() int {
	return e.intVal
}

// StringVal returns the string value, if it was set. Default "".
func (e *ExpressionResult) StringVal() string {
	return e.strVal
}

// RollHistory is a list of all roll results in string format.
//
// Mostly thread safe.
type RollHistory struct {
	rollResults []string
	accessMu    sync.Mutex
}

// ClearRolls clears the current roll history.
func (h *RollHistory) ClearRolls() {
	h.accessMu.Lock()
	defer h.accessMu.Unlock()
	h.rollResults = make([]string, 0)
}

// GetRollHistory returns a slice copy of the complete roll history.
func (h *RollHistory) GetRollHistory() []string {
	h.accessMu.Lock()
	defer h.accessMu.Unlock()
	return h.rollResults[:]
}

// AddRollToHistory adds the give roll string result to the history list.
func (h *RollHistory) AddRollToHistory(roll string) {
	h.accessMu.Lock()
	defer h.accessMu.Unlock()
	h.rollResults = append(h.rollResults, roll)
}

// LatestRoll returns the value of the last stored roll.
func (h *RollHistory) LatestRoll() string {
	h.accessMu.Lock()
	defer h.accessMu.Unlock()
	if len(h.rollResults) == 0 {
		return ""
	}
	return h.rollResults[len(h.rollResults)-1]
}

// ExecutionContext is a runtime context for scoping variable values,
// keepina consistent random number generator and referencing other tables.
type ExecutionContext struct {
	*RollHistory
	parent *ExecutionContext
	values map[string]*ExpressionResult
	packs  TableMap
	rand   RandomSource
}

// NewRootExecutionContext creates an empty ExecutionContext that is ready
// to be used, with empty history and the default random generator.
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

// SetRandom sets the RandomSource for this context.
func (ctx *ExecutionContext) SetRandom(r RandomSource) *ExecutionContext {
	ctx.rand = r
	return ctx
}

// Rand is a convenience method to get a random number from the context.
func (ctx *ExecutionContext) Rand(low int, high int) int {
	return ctx.rand.Get(low, high)
}

// Child creates a child execution context that can read-access the parent's variable values
// but writing to the same variables will not write those values to the parent definitions.
// This allows for recursion-based iteration.
//
// All other internal fields are direct copied to limit stack overflow issues and lightly optimise
// for fast access over traversing to the root node each evaluation.
func (ctx *ExecutionContext) Child() *ExecutionContext {
	return &ExecutionContext{
		RollHistory: ctx.RollHistory,
		parent:      ctx,
		values:      make(map[string]*ExpressionResult),
		packs:       ctx.packs,
		rand:        ctx.rand,
	}
}

// SetPacks assigns the table packs for this context.
func (ctx *ExecutionContext) SetPacks(packs TableMap) {
	ctx.packs = packs
}

// Set will set a variable named `key` to the value given.
func (ctx *ExecutionContext) Set(key string, val *ExpressionResult) {
	ctx.values[key] = val
}

// Resolve fetches key from the given context or the closest parent that has it
// defined. Nil is returned if the
func (ctx *ExecutionContext) Resolve(key string) (*ExpressionResult, error) {
	for c := ctx; c != nil; c = c.parent {
		if v, ok := c.values[key]; ok {
			return v, nil
		}
	}
	return nil, fmt.Errorf("variable accessed and not set: %s", key)
}

// EvaluateExpression Evaluates an expression with the given execution context outside
// of a Program. Mostly useful externally for compiler testing.
//
// This method contains the main evaluation loop that uses a slice for a program
// stack to prevent failing from deep call stacks.
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
	return nil, fmt.Errorf("this shouldn't happen, you have entered the matrix, have a fresh cookie: 0")
}

// TableMap is a type alias for mapping file hash keys to table definitions.
type TableMap map[string]*TablePack
