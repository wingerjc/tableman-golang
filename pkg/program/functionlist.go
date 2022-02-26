package program

var (
	genericFunctionList = map[string]*FunctionDef{
		"add": {
			funcName:    "add",
			minParams:   1,
			maxParams:   -1,
			resolve:     addResolve,
			verifyParam: onlyIntVerify,
		},
		"sum": {
			funcName:    "sum",
			minParams:   1,
			maxParams:   -1,
			resolve:     addResolve,
			verifyParam: onlyIntVerify,
		},
		"sub": {
			funcName:    "sub",
			minParams:   1,
			maxParams:   -1,
			resolve:     subResolve,
			verifyParam: onlyIntVerify,
		},
		"concat": {
			funcName:    "concat",
			minParams:   1,
			maxParams:   -1,
			resolve:     concatResolve,
			verifyParam: onlyStringVerify,
		},
		"upper": {
			funcName:    "upper",
			minParams:   1,
			maxParams:   1,
			resolve:     upperResolve,
			verifyParam: onlyStringVerify,
		},
		"lower": {
			funcName:    "lower",
			minParams:   1,
			maxParams:   1,
			resolve:     lowerResolve,
			verifyParam: onlyStringVerify,
		},
		"str": {
			funcName:    "str",
			minParams:   1,
			maxParams:   1,
			resolve:     toStrResolve,
			verifyParam: anyVerify,
		},
		"int": {
			funcName:    "int",
			minParams:   1,
			maxParams:   1,
			resolve:     toIntResolve,
			verifyParam: anyVerify,
		},
		"eq": {
			funcName:    "eq",
			minParams:   2,
			maxParams:   2,
			resolve:     eqResolve,
			verifyParam: anyVerify,
		},
		"gt": {
			funcName:    "gt",
			minParams:   2,
			maxParams:   2,
			resolve:     gtResolve,
			verifyParam: anyVerify,
		},
		"gte": {
			funcName:    "gte",
			minParams:   2,
			maxParams:   2,
			resolve:     gteResolve,
			verifyParam: anyVerify,
		},
		"lt": {
			funcName:    "lt",
			minParams:   2,
			maxParams:   2,
			resolve:     ltResolve,
			verifyParam: anyVerify,
		},
		"lte": {
			funcName:    "lte",
			minParams:   2,
			maxParams:   2,
			resolve:     lteResolve,
			verifyParam: anyVerify,
		},
		"and": {
			funcName:    "and",
			minParams:   2,
			maxParams:   -1,
			resolve:     andResolve,
			verifyParam: onlyIntVerify,
		},
		"or": {
			funcName:    "or",
			minParams:   2,
			maxParams:   -1,
			resolve:     orResolve,
			verifyParam: onlyIntVerify,
		},
		"not": {
			funcName:    "not",
			minParams:   1,
			maxParams:   1,
			resolve:     notResolve,
			verifyParam: onlyIntVerify,
		},
	}
	specializedFunctionList = map[string]func(string, []Evallable) (Evallable, error){
		"if": newIfFunction,
	}
)
