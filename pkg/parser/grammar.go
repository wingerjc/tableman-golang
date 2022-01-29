package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type TABLE_FILE struct {
	Header *FileHeader `@@?`
	Tables []*Table    `@@*`
}

type FileHeader struct {
	PackageName *string `@String`
}

type Table struct {
	Header *TableHeader `@@`
	Rows   []*TableRow  `@@+`
}

type TableHeader struct {
	Name *string `@Ident`
}

type TableRow struct {
	Default bool       `(@Default?`
	Numbers *RangeList `@@?`
	Label   *RowLabel  `@@? ":")?`
	Values  []*RowItem `@@+`
}

type RowItem struct {
	StringVal  *string     `@String`
	Expression *Expression `| @@`
}

type RowLabel struct {
	Single  *string `@TableName`
	Escaped *string `| @String`
}

type RangeList struct {
	First *NumberRange   `@@`
	Rest  []*NumberRange `(","@@)*`
}

// NumberRange represents a single number ror a range.
//
type NumberRange struct {
	First  *int `((@Number"-"`
	Last   int  `@Number)`
	Single int  `| @Number)`
}

type Expression struct {
	Value string `"{" @Number "}"`
}

type Roll struct {
}

const (
	NATURAL_NUMBER = `[1-9][0-9]*`
	WHOLE_NUMBER   = `0|([1-9][0-9]*)`
	INTEGER        = `0|(-?[1-9][0-9]*)`
)

var (
	fileLexer = lexer.MustStateful(lexer.Rules{
		"Root": []lexer.Rule{
			{"Default", `d(ef(ault)?)?`, nil},
			lexer.Include("Atomic"),
			{"Int", `[1-9][0-9]*`, nil},
			{"TableCall", `[Tt]\(`, lexer.Push("TableCall")},
			{"Expr", `{`, lexer.Push("Expr")},
			{"EOL", `[\n\r]+`, nil},
			{"ExtendLine", `->`, nil},
			{"Comment", `#.*$`, nil},
			{"Whitespace", `[ \t]+`, nil},
		},
		"Atomic": []lexer.Rule{
			{"TableName", `[a-zA-Z][a-zA-Z0-9_\-]*`, nil},
			{"Roll", `[1-9][0-9]*d[1-9][0-9]*((l|h)[1-9][0-9]+)?(\.(min|max|sum|avg|mode))?`, nil},
			{"String", `"(\\"|[^"])*"`, nil},
			{"Number", `[1-9][0-9]*`, nil},
			{"RangeDash", `-`, nil},
			{"ListDelimeter", `,`, nil},
			{"TableDelimeter", `:`, nil},
		},

		"Expr": []lexer.Rule{
			{"Number", NATURAL_NUMBER, nil},
			{"ExprEnd", `}`, lexer.Pop()},
		},
		"TableCall": []lexer.Rule{
			{"TableCallEnd", `\)`, lexer.Pop()},
		},
	})
)

func GetParser() (*participle.Parser, error) {
	return participle.Build(
		&TABLE_FILE{},
		participle.Lexer(fileLexer),
		participle.Elide("Comment", "Whitespace"),
	)
}
