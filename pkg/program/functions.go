package program

import (
	"fmt"
)

var (
	FUNCTION_LIST = map[string]*FunctionDef{
		"add": {funcName: "add", resolve: addResolve, verifyParam: addVerify},
		"sum": {funcName: "sum", resolve: addResolve, verifyParam: addVerify},
	}
)

type FunctionDef struct {
	funcName    string
	resolve     func([]*ExpressionResult) (*ExpressionResult, error)
	verifyParam func(ResultType, int) bool
}

type GenericFunction struct {
	params []Evallable
	config *FunctionDef
}

func NewFunction(name string, params []Evallable) (Evallable, error) {
	config, ok := FUNCTION_LIST[name]
	if !ok {
		return nil, fmt.Errorf("")
	}
	return &GenericFunction{
		params: params,
		config: config,
	}, nil
}

func (g *GenericFunction) Eval() ExpressionEval {
	return &EvalGenericFunc{
		funcDef: g,
		vals:    make([]*ExpressionResult, len(g.params)),
		index:   0,
	}
}

type EvalGenericFunc struct {
	ctx     *ExecutionContext
	funcDef *GenericFunction
	vals    []*ExpressionResult
	index   int
}

func (g *EvalGenericFunc) SetContext(ctx *ExecutionContext) ExpressionEval {
	g.ctx = ctx
	return g
}
func (g *EvalGenericFunc) HasNext() bool {
	return g.index < len(g.vals)
}

func (g *EvalGenericFunc) Next() ExpressionEval {
	res := g.funcDef.params[g.index]
	return res.Eval().SetContext(g.ctx.Child())
}

func (g *EvalGenericFunc) Provide(res *ExpressionResult) error {
	if !g.funcDef.config.verifyParam(res.resultType, g.index) {
		return fmt.Errorf("could not execute %s, wrong type for parameter %d", g.funcDef.config.funcName, g.index+1)
	}
	g.vals[g.index] = res
	g.index++
	return nil
}

func (g *EvalGenericFunc) Resolve() (*ExpressionResult, error) {
	return g.funcDef.config.resolve(g.vals)
}

func addResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	sum := 0
	for i, x := range results {
		if x.resultType != INT_RESULT {
			return nil, fmt.Errorf(
				"illegal type for 'add' at parameter %d, should be a number, got %s",
				i,
				x.StringVal(),
			)
		}
		sum += x.IntVal()
	}
	return NewIntResult(sum), nil
}

func addVerify(t ResultType, index int) bool {
	return t == INT_RESULT
}
