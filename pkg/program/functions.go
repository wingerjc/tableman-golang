package program

import (
	"fmt"
	"strconv"
	"strings"
)

type FunctionDef struct {
	funcName    string
	minParams   int
	maxParams   int
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
		return nil, fmt.Errorf("could not find function '%s'", name)
	}
	if config.minParams > len(params) {
		return nil, fmt.Errorf("too few params (%d) for function '%s', expected at least %d",
			len(params),
			config.funcName,
			config.minParams,
		)
	}
	if config.maxParams >= 0 && config.maxParams < len(params) {
		return nil, fmt.Errorf("too many params (%d) for function '%s'",
			len(params),
			config.funcName,
		)
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
		return fmt.Errorf("could not execute %s, wrong type for parameter %d in function '%s'",
			g.funcDef.config.funcName,
			g.index+1,
			g.funcDef.config.funcName,
		)
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
	for _, x := range results {
		sum += x.IntVal()
	}
	return NewIntResult(sum), nil
}

func subResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	total := results[0].IntVal()
	for _, r := range results[1:] {
		total -= r.IntVal()
	}
	return NewIntResult(total), nil
}

func concatResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	final := ""
	for _, r := range results {
		final += r.StringVal()
	}
	return NewStringResult(final), nil
}

func upperResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	return NewStringResult(strings.ToUpper(results[0].StringVal())), nil
}

func lowerResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	return NewStringResult(strings.ToLower(results[0].StringVal())), nil
}

func toStrResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	if results[0].MatchType(STRING_RESULT) {
		return results[0], nil
	}
	return NewStringResult(strconv.Itoa(results[0].IntVal())), nil
}

func toIntResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	if results[0].MatchType(INT_RESULT) {
		return results[0], nil
	}
	result, err := strconv.Atoi(results[0].StringVal())
	return NewIntResult(result), err
}

func onlyIntVerify(t ResultType, index int) bool {
	return t == INT_RESULT
}

func onlyStringVerify(t ResultType, index int) bool {
	return t == STRING_RESULT
}

func anyVerify(t ResultType, index int) bool {
	return true
}
