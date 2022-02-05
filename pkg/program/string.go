package program

type String struct {
	value   string
	isLabel bool
}

func NewString(value string, isLabel bool) Evallable {
	return &String{
		value:   value,
		isLabel: isLabel,
	}
}

func (s *String) IsLabel() bool {
	return s.isLabel
}

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

func (s *stringEval) Next() ExpressionEval {
	return nil
}

func (s *stringEval) Provide(res *ExpressionResult) error {
	return nil
}

func (s *stringEval) Resolve() (*ExpressionResult, error) {
	return NewStringResult(s.value), nil
}
