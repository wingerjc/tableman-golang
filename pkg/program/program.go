package program

import (
	"fmt"

	"github.com/wingerjc/tableman-golang/pkg/parser"
)

type Evallable interface {
	Eval() ExpressionEval
}

type ExpressionEval interface {
	SetContext(ctx *ExecutionContext) ExpressionEval
	HasNext() bool
	Next() ExpressionEval
	Provide(res *ExpressionResult) error
	Resolve() (*ExpressionResult, error)
}

type Program struct {
	packs map[string]*TablePack
}

func NewProgram() *Program {
	packs := make(map[string]*TablePack)
	return &Program{
		packs: packs,
	}
}

func (p *Program) AddPack(pack *TablePack) error {
	if _, ok := p.packs[pack.Name]; ok {
		return fmt.Errorf("pack name conflict, cannot have 2 packs named: %s", pack.Name)
	}
	p.packs[pack.Name] = pack
	return nil
}

type TablePack struct {
	Name    string
	tables  map[string]Table
	Imports []*Import
}

func TablePackFromAST(ast *parser.TableFile) (*TablePack, []error) {
	return &TablePack{
		tables:  make(map[string]Table),
		Imports: make([]*Import, 0),
	}, nil
}

type Import struct {
	FileName string
	Alias    string
}

type Table struct {
	// rows []*TableRow
}

type TableRow struct {
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
	if e.resultType == INT_RESULT {
		return e.intVal == other.intVal
	}
	if e.resultType == STRING_RESULT {
		return e.strVal == other.strVal
	}
	return false
}

func (e *ExpressionResult) SameType(other *ExpressionResult) bool {
	return e.resultType == other.resultType
}

func NewStringResult(val string) *ExpressionResult {
	return &ExpressionResult{
		resultType: STRING_RESULT,
		strVal:     val,
	}
}

func NewIntResult(val int) *ExpressionResult {
	return &ExpressionResult{
		resultType: INT_RESULT,
		intVal:     val,
	}
}

func (r *ExpressionResult) MatchType(types ...ResultType) bool {
	for _, t := range types {
		if r.resultType == t {
			return true
		}
	}
	return false
}

func (r *ExpressionResult) IntVal() int {
	return r.intVal
}

func (r *ExpressionResult) StringVal() string {
	return r.strVal
}

const (
	ANY_TYPE      = 0
	STRING_RESULT = 1
	INT_RESULT    = 2
)

type ExecutionContext struct {
	parent *ExecutionContext
	values map[string]*ExpressionResult
}

func NewRootExecutionContext() *ExecutionContext {
	return &ExecutionContext{
		parent: nil,
		values: make(map[string]*ExpressionResult),
	}
}

func (ctx *ExecutionContext) Child() *ExecutionContext {
	return &ExecutionContext{
		parent: ctx,
		values: make(map[string]*ExpressionResult),
	}
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

func EvaluateExpression(e Evallable) (*ExpressionResult, error) {
	stack := make([]ExpressionEval, 0)
	stack = append(stack, e.Eval())
	stack[0].SetContext(NewRootExecutionContext())
	for len(stack) > 0 {
		// See if we need to push another resolution node on the current stack.
		cur := stack[len(stack)-1]
		if cur.HasNext() {
			stack = append(stack, cur.Next())
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
