package program

type Expression struct {
	vars map[string]Evallable
	expr Evallable
}

func NewExpression(vars map[string]Evallable, expr Evallable) Evallable {
	return &Expression{
		vars: vars,
		expr: expr,
	}
}

func (e *Expression) Eval() ExpressionEval {
	keys := make([]string, 0, len(e.vars))
	for k := range e.vars {
		keys = append(keys, k)
	}
	return &runtimeExpression{
		expr:  e,
		keys:  keys,
		index: 0,
	}
}

type runtimeExpression struct {
	expr  *Expression
	ctx   *ExecutionContext
	res   *ExpressionResult
	keys  []string
	index int
}

func (r *runtimeExpression) SetContext(ctx *ExecutionContext) ExpressionEval {
	r.ctx = ctx
	return r
}

func (r *runtimeExpression) HasNext() bool {
	return r.index <= len(r.keys)
}

func (r *runtimeExpression) Next() ExpressionEval {
	if r.index < len(r.keys) {
		return r.expr.vars[r.currentKey()].Eval().SetContext(r.ctx.Child())
	}
	return r.expr.expr.Eval().SetContext(r.ctx.Child())
}

func (r *runtimeExpression) currentKey() string {
	return r.keys[r.index]
}

func (r *runtimeExpression) Provide(res *ExpressionResult) error {
	if r.index < len(r.keys) {
		r.ctx.Set(r.currentKey(), res)
		r.index++
		return nil
	}
	r.res = res
	r.index++
	return nil
}

func (r *runtimeExpression) Resolve() (*ExpressionResult, error) {
	return r.res, nil
}