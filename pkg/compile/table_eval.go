package compiler

import (
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

func compileTable(t *parser.Table, packKeys nameMap) (*program.Table, error) {
	tags := make(map[string]string)
	for _, tag := range t.Header.Tags {
		tags[tag.Key.String()] = tag.Value.String()
	}
	rows := make([]*program.TableRow, 0)
	if len(t.Rows) == 0 {
		numSteps := len(t.Generator.Steps)
		counts := make([]int, len(t.Generator.Steps))
		rangeInt := 1
		// Permute without recursion
		for {
			val := ""
			for i := 0; i < numSteps; i++ {
				val += t.Generator.Steps[i].StrVal(counts[i])
			}
			rows = append(rows, stringRow(val, rangeInt))
			var i int
			for i = 0; i < numSteps; i++ {
				counts[i]++
				if counts[i] == len(t.Generator.Steps[i].Values) {
					counts[i] = 0
				} else {
					break
				}
			}
			rangeInt++
			if i == numSteps {
				break
			}
		}
	} else {
		for _, r := range t.Rows {
			newRow, err := compileRow(r, packKeys)
			if err != nil {
				return nil, err
			}
			rows = append(rows, newRow)
		}
	}
	return program.NewTable(t.Header.Name, tags, rows), nil
}

func stringRow(val string, rangeInt int) *program.TableRow {
	rangeVal := make([]*program.Range, 1)
	rangeVal[0] = program.NewRange(rangeInt, rangeInt)
	return program.NewTableRow(
		"",
		rangeVal,
		1,
		1,
		false,
		program.NewString(val, false),
	)
}

func compileRow(r *parser.TableRow, packKeys nameMap) (*program.TableRow, error) {
	var err error
	items := make([]program.Evallable, 0)
	for _, i := range r.Values {
		var e program.Evallable
		if i.Expression != nil {
			e, err = compileExpression(i.Expression, packKeys)
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
