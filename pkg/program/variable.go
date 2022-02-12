package program

import "fmt"

type Variable struct {
	name string
}

func NewVariable(name string) Evallable {
	return &Variable{
		name: name,
	}
}

func (v *Variable) Eval() ExpressionEval {
	return &variableEval{
		name: v.name,
	}
}

type variableEval struct {
	name string
	ctx  *ExecutionContext
}

func (v *variableEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	v.ctx = ctx
	return v
}

func (v *variableEval) HasNext() bool {
	return false
}

func (v *variableEval) Next() (ExpressionEval, error) {
	return nil, fmt.Errorf("variables have no sub-expressions")
}

func (v *variableEval) Provide(res *ExpressionResult) error {
	return fmt.Errorf("variable type should never take results")
}

func (v *variableEval) Resolve() (*ExpressionResult, error) {
	res := v.ctx.Resolve(v.name)
	if res == nil {
		return nil, fmt.Errorf("variable not set: %s", v.name)
	}

	return res, nil
}
