package program

import "fmt"

// Variable is an evallable for a variable access.
type Variable struct {
	name string
}

// NewVariable creates a new variable access Evallable.
func NewVariable(name string) Evallable {
	return &Variable{
		name: name,
	}
}

// Eval implementation for Evallable interface.
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
