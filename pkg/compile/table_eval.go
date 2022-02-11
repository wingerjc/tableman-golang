package compiler

import (
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func CompileTable(t *parser.Table, random program.RandomSource, packKeys nameMap) (*program.Table, error) {
	tags := make(map[string]string)
	for _, tag := range t.Header.Tags {
		tags[tag.Key.String()] = tag.Value.String()
	}
	rows := make([]*program.TableRow, 0)
	for _, r := range t.Rows {
		newRow, err := CompileRow(r, packKeys)
		if err != nil {
			return nil, err
		}
		rows = append(rows, newRow)
	}
	return program.NewTable(t.Header.Name, tags, rows, random), nil
}

func CompileRow(r *parser.TableRow, packKeys nameMap) (*program.TableRow, error) {
	var err error
	items := make([]program.Evallable, 0)
	for _, i := range r.Values {
		var e program.Evallable
		if i.Expression != nil {
			e, err = CompileExpression(i.Expression, packKeys)
		} else {
			e = program.NewString(i.String(), false)
		}
		if err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	value := program.NewListExpression(items)
	label := ""
	if r.Label != nil {
		label = r.Label.String()
	}
	weight := 1
	if r.Weight > 1 {
		weight = r.Weight
	}
	count := 1
	if r.Count > 1 {
		count = r.Count
	}
	rangeVal := make([]*program.Range, 0)
	if r.Numbers != nil {
		for _, x := range r.Numbers.Ranges {
			if x.First != nil {
				rangeVal = append(rangeVal, program.NewRange(*x.First, x.Last))
			} else {
				rangeVal = append(rangeVal, program.NewRange(x.Single, x.Single))
			}
		}
	}
	return program.NewTableRow(label, rangeVal, weight, count, r.Default, value), nil
}
