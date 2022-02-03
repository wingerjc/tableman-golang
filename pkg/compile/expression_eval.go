package compiler

import (
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func CompileValueExpr(node *parser.ValueExpr) (program.Evallable, error) {
	switch node.GetType() {
	case parser.FUNC_EXPR_T:
		params, err := getParams(node)
		if err != nil {
			return nil, err
		}
		return program.NewFunction(node.Call.Name.Names[0], params)
	case parser.NUM_EXPR_T:
		return program.NewNumber(*node.Num), nil
	}
	return nil, nil
}

func getParams(node *parser.ValueExpr) ([]program.Evallable, error) {
	res := make([]program.Evallable, 0, len(node.Call.Params))
	for _, x := range node.Call.Params {
		expr, err := CompileValueExpr(x)
		if err != nil {
			return nil, err
		}
		res = append(res, expr)
	}
	return res, nil
}
