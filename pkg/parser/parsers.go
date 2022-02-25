package parser

import (
	"github.com/alecthomas/participle/v2"
)

// TableFileParser can parse full table files from text to AST.
type TableFileParser struct {
	p *participle.Parser
}

// Parse parses a string formatted as a table file to AST.
func (t *TableFileParser) Parse(code string) (*TableFile, error) {
	res := &TableFile{}
	err := t.p.ParseString("", code, res)
	return res, err
}

// GetParser creates a new parser for table files.
func GetParser() (*TableFileParser, error) {
	p, err := participle.Build(&TableFile{}, DefaultLexer, DefaultElide)
	if err != nil {
		return nil, err
	}
	return &TableFileParser{
		p: p,
	}, nil
}

// ExpressionParser can parse expressions from text to AST.
type ExpressionParser struct {
	p *participle.Parser
}

// Parse parses a string formatted as an expression to AST.
func (e *ExpressionParser) Parse(code string) (*Expression, error) {
	res := &Expression{}
	err := e.p.ParseString("", code, res)
	return res, err
}

// GetExpressionParser creates a new parser for expressions.
func GetExpressionParser() (*ExpressionParser, error) {
	p, err := participle.Build(&Expression{}, DefaultLexer, DefaultElide)
	if err != nil {
		return nil, err
	}
	return &ExpressionParser{
		p: p,
	}, nil
}

// RollParser can parse roll expressions from text to AST.
type RollParser struct {
	p *participle.Parser
}

// Parse parses a string formatted as a roll to AST.
func (r *RollParser) Parse(code string) (*Roll, error) {
	res := &Roll{}
	err := r.p.ParseString("", code, res)
	return res, err
}

// GetRollParser creates a new parser for rolls.
func GetRollParser() (*RollParser, error) {
	p, err := participle.Build(&Roll{}, DefaultLexer, DefaultElide)
	if err != nil {
		return nil, err
	}
	return &RollParser{
		p: p,
	}, nil
}

// RowParser can parse table row expressions from text to AST.
type RowParser struct {
	p *participle.Parser
}

// Parse parses a string formatted as a table row to AST.
func (e *RowParser) Parse(code string) (*TableRow, error) {
	res := &TableRow{}
	err := e.p.ParseString("", code, res)
	return res, err
}

// GetRowParser creates a new parser for table rows.
func GetRowParser() (*RowParser, error) {
	p, err := participle.Build(&TableRow{}, DefaultLexer, DefaultElide)
	if err != nil {
		return nil, err
	}
	return &RowParser{
		p: p,
	}, nil
}

// TableParser can parse table expressions from text to AST.
type TableParser struct {
	p *participle.Parser
}

// Parse parses a string formatted as a table to AST.
func (t *TableParser) Parse(code string) (*Table, error) {
	result := &Table{}
	err := t.p.ParseString("", code, result)
	// pp.Println(result)
	return result, err
}

// GetTableParser creates a new parse for tables.
func GetTableParser() (*TableParser, error) {
	p, err := participle.Build(&Table{}, DefaultLexer, DefaultElide)
	if err != nil {
		return nil, err
	}
	return &TableParser{
		p: p,
	}, nil
}
