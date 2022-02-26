package program

import "fmt"

func eqResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	val := 0
	if results[0].Equal(results[1]) {
		val = 1
	}
	return NewIntResult(val), nil
}

func gtResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	a := results[0]
	b := results[1]
	if !a.SameType(b) {
		return nil, fmt.Errorf("types do not match for function: %s", "gt")
	}
	val := 0
	if a.MatchType(IntResult) && a.IntVal() > b.IntVal() {
		val = 1
	}
	if a.MatchType(StringResult) && a.StringVal() > b.StringVal() {
		val = 1
	}
	return NewIntResult(val), nil
}

func gteResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	a := results[0]
	b := results[1]
	if !a.SameType(b) {
		return nil, fmt.Errorf("types do not match for function: %s", "te")
	}
	val := 0
	if a.MatchType(IntResult) && a.IntVal() >= b.IntVal() {
		val = 1
	}
	if a.MatchType(StringResult) && a.StringVal() >= b.StringVal() {
		val = 1
	}
	return NewIntResult(val), nil
}

func ltResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	a := results[0]
	b := results[1]
	if !a.SameType(b) {
		return nil, fmt.Errorf("types do not match for function: %s", "lt")
	}
	val := 0
	if a.MatchType(IntResult) && a.IntVal() < b.IntVal() {
		val = 1
	}
	if a.MatchType(StringResult) && a.StringVal() < b.StringVal() {
		val = 1
	}
	return NewIntResult(val), nil
}

func lteResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	a := results[0]
	b := results[1]
	if !a.SameType(b) {
		return nil, fmt.Errorf("types do not match for function: %s", "lte")
	}
	val := 0
	if a.MatchType(IntResult) && a.IntVal() <= b.IntVal() {
		val = 1
	}
	if a.MatchType(StringResult) && a.StringVal() <= b.StringVal() {
		val = 1
	}
	return NewIntResult(val), nil
}

func orResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	for _, v := range results {
		if v.IntVal() != 0 {
			return NewIntResult(1), nil
		}
	}
	return NewIntResult(0), nil
}

func andResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	for _, v := range results {
		if v.IntVal() == 0 {
			return NewIntResult(0), nil
		}
	}
	return NewIntResult(1), nil
}

func notResolve(results []*ExpressionResult) (*ExpressionResult, error) {
	if results[0].IntVal() == 0 {
		return NewIntResult(1), nil
	}
	return NewIntResult(0), nil
}
