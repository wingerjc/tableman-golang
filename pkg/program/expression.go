package program

import "github.com/k0kubun/pp"

type Expression struct {
	varOrder []string
	vars     map[string]Evallable
	expr     Evallable
}

func NewExpression(varOrder []string, vars map[string]Evallable, expr Evallable) Evallable {
	return &Expression{
		varOrder: varOrder,
		vars:     vars,
		expr:     expr,
	}
}

func (e *Expression) Eval() ExpressionEval {
	return &runtimeExpression{
		expr:  e,
		keys:  e.varOrder,
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
	pp.Println(r.ctx)
	return r.res, nil
}
