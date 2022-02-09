package program

type TableCall struct {
}

func NewTableCall(packageKey string, tableName string, params []Evallable) Evallable {
	return &TableCall{}
}

func (c *TableCall) Eval() ExpressionEval {
	return &tableCallEval{}
}

type tableCallEval struct {
}

func (t *tableCallEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	return t
}

func (t *tableCallEval) HasNext() bool {
	return false
}

func (t *tableCallEval) Next() ExpressionEval {
	return nil
}

func (t *tableCallEval) Provide(res *ExpressionResult) error {
	return nil
}

func (t *tableCallEval) Resolve() (*ExpressionResult, error) {
	return nil, nil
}
