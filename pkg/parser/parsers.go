package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/k0kubun/pp"
)

type TableFileParser struct {
	p *participle.Parser
}

func (t *TableFileParser) Parse(code string) (*TableFile, error) {
	res := &TableFile{}
	err := t.p.ParseString("", code, res)
	return res, err
}

func GetParser() (*TableFileParser, error) {
	p, err := participle.Build(&TableFile{}, DEFAULT_LEXER, DEFAULT_ELIDE)
	if err != nil {
		return nil, err
	}
	return &TableFileParser{
		p: p,
	}, nil
}

type ExpressionParser struct {
	p *participle.Parser
}

func (e *ExpressionParser) Parse(code string) (*Expression, error) {
	res := &Expression{}
	err := e.p.ParseString("", code, res)
	return res, err
}

func GetExpressionParser() (*ExpressionParser, error) {
	p, err := participle.Build(&Expression{}, DEFAULT_LEXER, DEFAULT_ELIDE)
	if err != nil {
		return nil, err
	}
	return &ExpressionParser{
		p: p,
	}, nil
}

type RowParser struct {
	p *participle.Parser
}

func (e *RowParser) Parse(code string) (*TableRow, error) {
	res := &TableRow{}
	err := e.p.ParseString("", code, res)
	return res, err
}

func GetRowParser() (*RowParser, error) {
	p, err := participle.Build(&TableRow{}, DEFAULT_LEXER, DEFAULT_ELIDE)
	if err != nil {
		return nil, err
	}
	return &RowParser{
		p: p,
	}, nil
}

type TableParser struct {
	p *participle.Parser
}

func (t *TableParser) Parse(code string) (*Table, error) {
	result := &Table{}
	err := t.p.ParseString("", code, result)
	pp.Println(result)
	return result, err
}

func GetTableParser() (*TableParser, error) {
	p, err := participle.Build(&Table{}, DEFAULT_LEXER, DEFAULT_ELIDE)
	if err != nil {
		return nil, err
	}
	return &TableParser{
		p: p,
	}, nil
}
