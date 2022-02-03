package program

type Number struct {
	value int
}

func NewNumber(value int) Evallable {
	return &Number{
		value: value,
	}
}

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

func (n *numberEval) Next() ExpressionEval {
	return nil
}

func (n *numberEval) Provide(res *ExpressionResult) error {
	return nil
}

func (n *numberEval) Resolve() (*ExpressionResult, error) {
	return NewIntResult(n.value), nil
}
