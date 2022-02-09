package compiler

import (
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func CompileExpression(node *parser.Expression) (program.Evallable, error) {
	vars := make(map[string]program.Evallable)
	varOrder := make([]string, 0, len(node.Vars))
	for _, v := range node.Vars {
		expr, err := CompileValueExpr(v.AssignedValue)
		if err != nil {
			return nil, err
		}
		vars[v.VarName.Name] = expr
		varOrder = append(varOrder, v.VarName.Name)
	}
	res, err := CompileValueExpr(node.Value)
	if err != nil {
		return nil, err
	}
	return program.NewExpression(varOrder, vars, res), nil
}

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
	case parser.LABEL_EXPR_T:
		return program.NewString(node.Label.String(), node.Label.IsLabel()), nil
	case parser.VAR_EXPR_T:
		return program.NewVariable(node.Variable.Name), nil
	case parser.TABLE_EXPR_T:
		params, err := getParams(node)
		if err != nil {
			return nil, err
		}
		return program.NewTableCall(node.Call.Name.PackageName(), node.Call.Name.TableName(), params), nil
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
