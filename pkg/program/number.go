package program

import "fmt"

// Number is an Evallable for a numeric constant.
type Number struct {
	value int
}

// NewNumber creates a new numeric Evallable.
func NewNumber(value int) Evallable {
	return &Number{
		value: value,
	}
}

// Eval implementation for Evallable interface.
func (n *Number) Eval() ExpressionEval {
	return newNumberEval(n.value)
}

type numberEval struct {
	value int
}

func newNumberEval(val int) ExpressionEval {
	return &numberEval{
		value: val,
	}
}

func (n *numberEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	return n
}

func (n *numberEval) HasNext() bool {
	return false
}

func (n *numberEval) Next() (ExpressionEval, error) {
	return nil, fmt.Errorf("number expressions should not have sub-expressions")
}

func (n *numberEval) Provide(res *ExpressionResult) error {
	return nil
}

func (n *numberEval) Resolve() (*ExpressionResult, error) {
	return NewIntResult(n.value), nil
}
