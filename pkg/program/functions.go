package program

import (
	"fmt"
	"strconv"
	"strings"
)

// FunctionDef is a way to define a function so it can be used
// in a modular way.
type FunctionDef struct {
	funcName    string
	minParams   int
	maxParams   int
	resolve     func([]*ExpressionResult) (*ExpressionResult, error)
	verifyParam func(ResultType, int) bool
}

// GenericFunction allows a simple FunctionDef to be wrapped for simpler definitions.
// An Evallable.
type GenericFunction struct {
	params []Evallable
	config *FunctionDef
}

// NewFunction creates a function Evallable, erroring if it can on the
// format of the parameters.
func NewFunction(name string, params []Evallable) (Evallable, error) {
	fn, ok := specializedFunctionList[name]
	if ok {
		return fn(name, params)
	}
	config, ok := genericFunctionList[name]
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

// Eval implementation for Evallable interface.
func (g *GenericFunction) Eval() ExpressionEval {
	return &evalGenericFunc{
		funcDef: g,
		vals:    make([]*ExpressionResult, len(g.params)),
		index:   0,
	}
}

type evalGenericFunc struct {
	ctx     *ExecutionContext
	funcDef *GenericFunction
	vals    []*ExpressionResult
	index   int
}

func (g *evalGenericFunc) SetContext(ctx *ExecutionContext) ExpressionEval {
	g.ctx = ctx
	return g
}
func (g *evalGenericFunc) HasNext() bool {
	return g.index < len(g.vals)
}

func (g *evalGenericFunc) Next() (ExpressionEval, error) {
	if g.index > len(g.funcDef.params) {
		return nil, fmt.Errorf("accessing too many sub-expressions for function call")
	}
	res := g.funcDef.params[g.index]
	return res.Eval().SetContext(g.ctx.Child()), nil
}

func (g *evalGenericFunc) Provide(res *ExpressionResult) error {
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

func (g *evalGenericFunc) Resolve() (*ExpressionResult, error) {
	return g.funcDef.config.resolve(g.vals)
}

func addResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	sum := 0
	for _, x := range results {
		sum += x.IntVal()
	}
	return NewIntResult(sum), nil
}

func multResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	product := 1
	for _, x := range results {
		product *= x.intVal
	}
	return NewIntResult(product), nil
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
	if results[0].MatchType(StringResult) {
		return results[0], nil
	}
	return NewStringResult(strconv.Itoa(results[0].IntVal())), nil
}

func toIntResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	if results[0].MatchType(IntResult) {
		return results[0], nil
	}
	result, err := strconv.Atoi(results[0].StringVal())
	return NewIntResult(result), err
}

func onlyIntVerify(t ResultType, index int) bool {
	return t == IntResult
}

func onlyStringVerify(t ResultType, index int) bool {
	return t == StringResult
}

func anyVerify(t ResultType, index int) bool {
	return true
}

type ifFunction struct {
	condition Evallable
	trueVal   Evallable
	falseVal  Evallable
}

func newIfFunction(name string, vals []Evallable) (Evallable, error) {
	if len(vals) != 3 {
		return nil, fmt.Errorf("need 3 parameters for 'if', was passed %d", len(vals))
	}
	return &ifFunction{
		condition: vals[0],
		trueVal:   vals[1],
		falseVal:  vals[2],
	}, nil
}

func (i *ifFunction) Eval() ExpressionEval {
	return &ifFunctionEval{
		config: i,
	}
}

type ifFunctionEval struct {
	ctx             *ExecutionContext
	config          *ifFunction
	conditionResult *ExpressionResult
	result          *ExpressionResult
}

func (i *ifFunctionEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	i.ctx = ctx
	return i
}

func (i *ifFunctionEval) HasNext() bool {
	return i.conditionResult == nil || i.result == nil
}

func (i *ifFunctionEval) Next() (ExpressionEval, error) {
	if i.conditionResult == nil {
		return i.config.condition.Eval().SetContext(i.ctx.Child()), nil
	}
	if i.conditionResult.BoolVal() {
		return i.config.trueVal.Eval().SetContext(i.ctx.Child()), nil
	}
	return i.config.falseVal.Eval().SetContext(i.ctx.Child()), nil
}

func (i *ifFunctionEval) Provide(res *ExpressionResult) error {
	if i.conditionResult == nil {
		if !res.MatchType(IntResult) {
			return fmt.Errorf("'if' condition must be an integer expression")
		}
		i.conditionResult = res
		return nil
	}
	if i.result != nil {
		return fmt.Errorf("'if' result already set, cannot set a second time")
	}
	i.result = res
	return nil
}

func (i *ifFunctionEval) Resolve() (*ExpressionResult, error) {
	return i.result, nil
}
