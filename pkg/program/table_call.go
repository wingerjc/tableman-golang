package program

import (
	"fmt"
)

// TableCall is an Evallable for calls to a table.
type TableCall struct {
	packageKey  string
	packageName string
	tableName   string
	params      []Evallable
}

// NewTableCall creates a new table call Evallable.
func NewTableCall(
	packageKey string,
	packageName string,
	tableName string,
	params []Evallable,
) (Evallable, error) {
	if len(packageKey) == 0 {
		packageKey = RootPack
	}
	if len(params) > 2 {
		return nil, fmt.Errorf("table call cannot have more than 2 parameters")
	}
	return &TableCall{
		packageKey:  packageKey,
		packageName: packageName,
		tableName:   tableName,
		params:      params,
	}, nil
}

// Eval implementation for Evallable interface.
func (c *TableCall) Eval() ExpressionEval {
	return &tableCallEval{
		def:        c,
		results:    make([]*ExpressionResult, 0, len(c.params)),
		index:      0,
		paramCount: len(c.params),
	}
}

type tableCallEval struct {
	ctx         *ExecutionContext
	def         *TableCall
	results     []*ExpressionResult
	tableResult *ExpressionResult
	paramCount  int
	index       int
}

func (t *tableCallEval) SetContext(ctx *ExecutionContext) ExpressionEval {
	t.ctx = ctx
	return t
}

func (t *tableCallEval) HasNext() bool {
	return t.index <= t.paramCount
}

func (t *tableCallEval) Next() (ExpressionEval, error) {
	if t.index == t.paramCount {
		return t.callTable()
	}
	return t.def.params[t.index].Eval().SetContext(t.ctx), nil
}

func (t *tableCallEval) callTable() (ExpressionEval, error) {
	pack, ok := t.ctx.packs[t.def.packageKey]
	if !ok {
		return nil, fmt.Errorf("could not access table pack '%s'", t.def.packageName)
	}
	table, ok := pack.tables[t.def.tableName]
	if !ok {
		return nil, fmt.Errorf("package '%s' has no table '%s'", t.def.packageName, t.def.tableName)
	}
	if t.paramCount == 0 {
		return table.Roll().Eval().SetContext(t.ctx.Child()), nil
	}
	if !t.results[0].MatchType(StringResult) {
		return nil, fmt.Errorf("roll type must be a string value one of: roll/weighted/index/label/deck")
	}
	switch t.results[0].StringVal() {
	case "roll":
		return table.Roll().Eval().SetContext(t.ctx.Child()), nil
	case "weighted":
		return table.WeightedRoll().Eval().SetContext(t.ctx.Child()), nil
	case "index":
		if t.paramCount != 2 {
			return nil, fmt.Errorf("index rolls require 2 parameters: '!t(index, <number>)'")
		}
		if !t.results[1].MatchType(IntResult) {
			return nil, fmt.Errorf("index rolls must have a number for the second parameter")
		}
		row, err := table.IndexRoll(t.results[1].IntVal())
		if err != nil {
			return nil, err
		}
		return row.Eval().SetContext(t.ctx.Child()), nil
	case "label":
		if t.paramCount != 2 {
			return nil, fmt.Errorf("label rolls require 2 parameters: '!t(label, <string>)")
		}
		if !t.results[1].MatchType(StringResult) {
			return nil, fmt.Errorf("label rolls must have a string value for the second parameter")
		}
		row, err := table.LabelRoll(t.results[1].StringVal())
		if err != nil {
			return nil, err
		}
		return row.Eval().SetContext(t.ctx.Child()), nil
	case "deck":
		if t.paramCount == 2 {
			shuf := t.results[1]
			if !shuf.MatchType(StringResult) || (shuf.StringVal() != "shuffle" && shuf.strVal != "no-shuffle") {
				return nil, fmt.Errorf("if passing 2 parameters to a deck roll the 2nd parameter should be 'shuffle' or 'no-shuffle'")
			}
		}
		if t.paramCount > 2 {
			return nil, fmt.Errorf("too many parameters for deck roll, should be 1 or 2")
		}
		if t.paramCount == 2 && t.results[1].strVal == "shuffle" {
			table.Shuffle()
		}
		row, err := table.DeckDraw()
		if err != nil {
			return nil, err
		}
		return row.Eval().SetContext(t.ctx.Child()), nil
	}
	return nil, fmt.Errorf("roll type must be a string value one of: roll/weighted/index/label/deck")
}

func (t *tableCallEval) Provide(res *ExpressionResult) error {
	if t.index == t.paramCount {
		t.tableResult = res
		t.index++
		return nil
	}
	if t.index > t.paramCount {
		return fmt.Errorf("too many sub-expression results applied to table call")
	}
	t.results = append(t.results, res)
	t.index++
	return nil
}

func (t *tableCallEval) Resolve() (*ExpressionResult, error) {
	if t.index == t.paramCount+1 {
		return t.tableResult, nil
	}
	return nil, fmt.Errorf("can't resolve table call, not all sub-expressions evaluated")
}
