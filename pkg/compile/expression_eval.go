package compiler

import (
	"fmt"

	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func CompileExpression(node *parser.Expression, packKeys nameMap) (program.Evallable, error) {
	vars := make(map[string]program.Evallable)
	varOrder := make([]string, 0, len(node.Vars))
	for _, v := range node.Vars {
		expr, err := CompileValueExpr(v.AssignedValue, packKeys)
		if err != nil {
			return nil, err
		}
		vars[v.VarName.Name] = expr
		varOrder = append(varOrder, v.VarName.Name)
	}
	res, err := CompileValueExpr(node.Value, packKeys)
	if err != nil {
		return nil, err
	}
	return program.NewExpression(varOrder, vars, res), nil
}

func CompileValueExpr(node *parser.ValueExpr, packKeys nameMap) (program.Evallable, error) {
	switch node.GetType() {
	case parser.FUNC_EXPR_T:
		params, err := getParams(node, packKeys)
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
		return compileTableCall(node, packKeys)
	case parser.ROLL_EXPR_T:
		return compileRollExpr(node.Roll)
	}
	return nil, fmt.Errorf("unkown expression type %s", node.GetStringType())
}

func compileRollExpr(node *parser.Roll) (program.Evallable, error) {
	count, sides, err := node.Dice()
	if err != nil {
		return nil, err
	}
	res := program.NewRoll(count, sides).
		WithPrint(node.Print).
		WithAggr(node.FnAggr())

	if len(node.RollCountAggrs) > 0 {
		aggrMap := make(map[int]*program.RollCountAggr)
		aggrList := make([]*program.RollCountAggr, 0)
		for _, a := range node.RollCountAggrs {
			if _, ok := aggrMap[a.Number]; ok {
				return nil, fmt.Errorf("double roll count aggrs assigned to number %d", a.Number)
			}
			r := program.NewRollCountAggr(a.Number, a.FinalMult())
			aggrMap[a.Number] = r
			aggrList = append(aggrList, r)
		}
		res = res.WithCountAggr(aggrList)
	}

	if len(node.RollSubset) > 0 {
		res = res.WithSelector(program.NewRollSelect(node.RollSubset == "h", node.SubsetCount))
	}

	return res, nil
}

func compileTableCall(node *parser.ValueExpr, packKeys nameMap) (program.Evallable, error) {
	params, err := getParams(node, packKeys)
	if err != nil {
		return nil, err
	}
	packName := node.Call.Name.PackageName()
	key, ok := packKeys[packName]
	if !ok {
		return nil, fmt.Errorf("could not find package '%s' did you forget or mistype an import?", packName)
	}
	return program.NewTableCall(
		key,
		packName,
		node.Call.Name.TableName(),
		params,
	)
}

func getParams(node *parser.ValueExpr, packKeys nameMap) ([]program.Evallable, error) {
	res := make([]program.Evallable, 0, len(node.Call.Params))
	for _, x := range node.Call.Params {
		expr, err := CompileValueExpr(x, packKeys)
		if err != nil {
			return nil, err
		}
		res = append(res, expr)
	}
	return res, nil
}
