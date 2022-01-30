package parser

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type TableFile struct {
	Pos    lexer.Position
	Header *FileHeader `parser:"EOL* @@"`
	Tables []*Table    `parser:"(EOL+ @@? TableBarrier?)*"`
}

type FileHeader struct {
	Pos             lexer.Position
	RootPackageName string             `parser:"PkgStart @TableName"`
	SubPackages     []string           `parser:"(PkgDelimiter @TableName)*"`
	Imports         []*ImportStatement `parser:"(EOL @@)*"`
}

type ImportStatement struct {
	Pos      lexer.Position
	FileName string "Import @FilePath"
}

type Table struct {
	Pos    lexer.Position
	Header *TableHeader `parser:"@@"`
	Rows   []*TableRow  `parser:"(EOL @@)+"`
}

type TableHeader struct {
	Pos  lexer.Position
	Name string `parser:"TableStart @TableName"`
	Tags []*Tag `parser:"(EOL+ @@)*"`
}

type Tag struct {
	Pos   lexer.Position
	Key   RowLabel `parser:"TagStart @@ TableDelimiter"`
	Value RowLabel `parser:"@@"`
}

type TableRow struct {
	Pos     lexer.Position
	Default bool       `parser:"(@Default?"`
	Weight  string     `parser:"@WeightMarker?"`
	Numbers *RangeList `parser:"@@?"`
	Label   *RowLabel  `parser:"@@? ':')?"`
	Values  []*RowItem `parser:"@@+"`
}

type RowItem struct {
	Pos        lexer.Position
	StringVal  *string     `parser:"(@String"`
	Expression *Expression `parser:"| @@)(ExtendLine EOL)?"`
}

type RowLabel struct {
	Pos     lexer.Position
	Single  *string `parser:"@TableName"`
	Escaped *string `parser:"| @String"`
}

type RangeList struct {
	Pos   lexer.Position
	First *NumberRange   `parser:"@@"`
	Rest  []*NumberRange `parser:"(','@@)*"`
}

// NumberRange represents a single number ror a range.
//
type NumberRange struct {
	Pos    lexer.Position
	First  *int `parser:"((@Number'-'"`
	Last   int  `parser:"@Number)"`
	Single int  `parser:"| @Number)"`
}

type Expression struct {
	Pos   lexer.Position
	Value string `parser:"'{' @Number '}'"`
}

type Roll struct {
	Pos            lexer.Position
	RollDice       string   `parser:"@Roll"`
	RollSubset     string   `parser:"@RollSubset? ("`
	RollFuncAggr   string   `parser:"@RollFuncAggr"`
	RollCountAggrs []string `parser:"|(@RollCountAggr+))?"`
}

const (
	NATURAL_NUMBER = `([1-9][0-9]*)`
	WHOLE_NUMBER   = `(0|([1-9][0-9]*))`
	INTEGER        = `(0|(-?[1-9][0-9]*))`
)

var (
	fileLexer = lexer.MustStateful(lexer.Rules{
		"Root": []lexer.Rule{
			{Name: "Default", Pattern: `d(ef(ault)?)?`},
			{Name: "PkgStart", Pattern: `TablePack`},
			{Name: "PkgDelimiter", Pattern: `\.`},
			{Name: "WeightMarker", Pattern: `w=` + WHOLE_NUMBER},
			{Name: "ExtendLine", Pattern: `->`},
			{Name: "TableBarrier", Pattern: `--(-+)`},
			{Name: "TableStart", Pattern: `[Tt](able)?:`},
			{Name: "Import", Pattern: `(?i)import`},
			{Name: "FilePath", Pattern: `(?i)f"(([A-Z]:)|~|(\.\.?))/.*"`},
			lexer.Include("Atomic"),
			{Name: "ListDelimeter", Pattern: `,`},
			{Name: "TableDelimiter", Pattern: `:`},
			{Name: "RangeDash", Pattern: `-`},
			{Name: "TableCall", Pattern: `[Tt]\(`, Action: lexer.Push("TableCall")},
			{Name: "Expr", Pattern: `{`, Action: lexer.Push("Expr")},
			{Name: "EOL", Pattern: `\r?\n`},
			{Name: "Comment", Pattern: `#.*$`},
			{Name: "CommentLine", Pattern: `^[ \t]*#.*\r?\n`},
			{Name: "Whitespace", Pattern: `[ \t]+`},
			{Name: "TagStart", Pattern: `~`},
		},
		"Atomic": []lexer.Rule{
			{Name: "TableName", Pattern: `[a-zA-Z][a-zA-Z0-9_\-]*`},
			{Name: "Roll", Pattern: NATURAL_NUMBER + `d` + NATURAL_NUMBER, Action: lexer.Push("Roll")},
			{Name: "String", Pattern: `"(\\"|[^"])*"`},
			{Name: "Number", Pattern: NATURAL_NUMBER},
		},
		"Roll": []lexer.Rule{
			{Name: "RollSubset", Pattern: `(l|h)` + NATURAL_NUMBER},
			{Name: "RollFuncAggr", Pattern: `\.(min|max|sum|avg|mode)`},
			{Name: "RollCountAggr", Pattern: `\.[+-]` + NATURAL_NUMBER + `(x` + NATURAL_NUMBER + `)?`},
			{Name: "RollEnd", Pattern: `[ \t]+|$|(\r?\n)`, Action: lexer.Pop()},
		},
		"Expr": []lexer.Rule{
			lexer.Include("Atomic"),
			{Name: "ExprEnd", Pattern: `}`, Action: lexer.Pop()},
		},
		"TableCall": []lexer.Rule{
			{Name: "TableCallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
	})
)

func GetParser() (*participle.Parser, error) {
	return participle.Build(
		&TableFile{},
		participle.Lexer(fileLexer),
		participle.Elide("Comment", "Whitespace", "CommentLine"),
	)
}
