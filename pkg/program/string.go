package program

import "fmt"

// String is an Evallable for a string value.
type String struct {
	value   string
	isLabel bool
}

// NewString creates a new string evallable.
func NewString(value string, isLabel bool) Evallable {
	return &String{
		value:   value,
		isLabel: isLabel,
	}
}

// IsLabel returns whether the string was considered a label at parse time.
func (s *String) IsLabel() bool {
	return s.isLabel
}

// Eval implementation for the Evallable inerface.
func (s *String) Eval() ExpressionEval {
	return &stringEval{
		value: s.value,
	}
}

type stringEval struct {
	value string
}

func (s *stringEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	return s
}

func (s *stringEval) HasNext() bool {
	return false
}

func (s *stringEval) Next() (ExpressionEval, error) {
	return nil, fmt.Errorf("string values should not have sub-expressions")
}

func (s *stringEval) Provide(res *ExpressionResult) error {
	return nil
}

func (s *stringEval) Resolve() (*ExpressionResult, error) {
	return NewStringResult(s.value), nil
}
