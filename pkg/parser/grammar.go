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
	Pos     lexer.Position
	Name    *ExtendedTableName `parser:"PkgStart @@"`
	Imports []*ImportStatement `parser:"(EOL @@)*"`
}

type ImportStatement struct {
	Pos      lexer.Position
	FileName string             `parser:"Import @FilePath"`
	Alias    *ExtendedTableName `parser:"(PackAlias @@)?"`
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
	Key   LabelString `parser:"TagStart @@ TableDelimiter"`
	Value LabelString `parser:"@@"`
}

type TableRow struct {
	Pos     lexer.Position
	Default bool         `parser:"(@Default?"`
	Weight  string       `parser:"@WeightMarker?"`
	Numbers *RangeList   `parser:"@@?"`
	Label   *LabelString `parser:"@@? ':')?"`
	Values  []*RowItem   `parser:"@@+"`
}

type RowItem struct {
	Pos        lexer.Position
	StringVal  *string     `parser:"(@String"`
	Expression *Expression `parser:"| @@)(ExtendLine EOL)?"`
}

type LabelString struct {
	Pos     lexer.Position
	Single  *string `parser:"@TableName"`
	Escaped *string `parser:"| @String"`
}

type RangeList struct {
	Pos   lexer.Position
	First *NumberRange   `parser:"@@"`
	Rest  []*NumberRange `parser:"(','@@)*"`
}

// NumberRange represents a single number or a range.
//
type NumberRange struct {
	Pos    lexer.Position
	First  *int `parser:"((@Number'-'"`
	Last   int  `parser:"@Number)"`
	Single int  `parser:"| @Number)"`
}

type Roll struct {
	Pos            lexer.Position
	RollDice       string   `parser:"@Roll"`
	RollSubset     string   `parser:"@RollSubset? ("`
	RollFuncAggr   string   `parser:"@RollFuncAggr"`
	RollCountAggrs []string `parser:"|(@RollCountAggr+))? RollEnd"`
}

type Expression struct {
	Pos   lexer.Position
	Vars  []*VariableDef `parser:"ExprStart (@@ (ListDelimiter @@)* EndVarList)?"`
	Value *ValueExpr     `parser:"@@ ExprEnd"`
}

type VariableDef struct {
	VarName       *VarName   `parser:"@@ VarAssign"`
	AssignedValue *ValueExpr `parser:"@@"`
}

type Call struct {
	IsTable bool              `parser:"@TableCallSignal?"`
	Name    ExtendedTableName `parser:"@@ CallStart"`
	Params  []*ValueExpr      `parser:"@@? (ListDelimiter @@)* CallEnd"`
}

type ValueExpr struct {
	Roll     *Roll        `parser:"@@"`
	Num      int          `parser:"| @Number"`
	IntVal   int          `parser:"| @Integer"`
	Call     *Call        `parser:"| @@"`
	Label    *LabelString `parser:"| @@"`
	Variable *VarName     `parser:"| @@"`
}

type VarName struct {
	Name string `parser:"VarPrefix @TableName"`
}
type ExtendedTableName struct {
	Names []string `parser:" @TableName (PkgDelimiter @TableName)*"`
}

const (
	NATURAL_NUMBER = `([1-9][0-9]*)`
	WHOLE_NUMBER   = `(0|([1-9][0-9]*))`
	INTEGER        = `(0|(-?[1-9][0-9]*))`
	IDENTIFIER     = `[a-zA-Z][a-zA-Z0-9\-_]*`
)

var (
	fileLexer = lexer.MustStateful(lexer.Rules{
		"Root": []lexer.Rule{
			lexer.Include("Whitespace"),
			{Name: "Default", Pattern: `Default`},
			{Name: "PkgStart", Pattern: `TablePack:`},
			{Name: "TableStart", Pattern: `TableDef:`},
			{Name: "Import", Pattern: `Import:`},
			{Name: "PackAlias", Pattern: `As:`},
			{Name: "WeightMarker", Pattern: `w=` + WHOLE_NUMBER},
			{Name: "ExtendLine", Pattern: `->`},
			{Name: "TableBarrier", Pattern: `--(-+)`},
			{Name: "FilePath", Pattern: `f\"(([A-Za-z]:)|~|(\.\.?))?/.*\"`},
			lexer.Include("Atomic"),
			{Name: "TableDelimiter", Pattern: `:`},
			{Name: "RangeDash", Pattern: `-`},
			{Name: "EOL", Pattern: `\r?\n`},
			{Name: "TagStart", Pattern: `~`},
		},
		"Atomic": []lexer.Rule{
			{Name: "TableName", Pattern: IDENTIFIER},
			{Name: "Roll", Pattern: NATURAL_NUMBER + `d` + NATURAL_NUMBER, Action: lexer.Push("Roll")},
			{Name: "CallStart", Pattern: `\(`, Action: lexer.Push("Call")},
			{Name: "ExprStart", Pattern: `{`, Action: lexer.Push("Expr")},
			{Name: "String", Pattern: `"(\\"|[^"])*"`},
			{Name: "Number", Pattern: NATURAL_NUMBER},
			{Name: "TableCallSignal", Pattern: `\!`},
			{Name: "VarPrefix", Pattern: `@`},
			{Name: "PkgDelimiter", Pattern: `\.`},
			{Name: "ListDelimiter", Pattern: `,`},
		},
		"Roll": []lexer.Rule{
			{Name: "RollSubset", Pattern: `(l|h)` + NATURAL_NUMBER},
			{Name: "RollFuncAggr", Pattern: `\.(min|max|sum|avg|mode|roll)`},
			{Name: "RollCountAggr", Pattern: `\.[+-]` + NATURAL_NUMBER + `(x` + NATURAL_NUMBER + `)?`},
			{Name: "RollEnd", Pattern: `\?`, Action: lexer.Pop()},
		},
		"Expr": []lexer.Rule{
			lexer.Include("Whitespace"),
			lexer.Include("Atomic"),
			lexer.Include("ExprValues"),
			{Name: "VarAssign", Pattern: `=`},
			{Name: "EndVarList", Pattern: `;`},
			{Name: "ExprEnd", Pattern: `\}`, Action: lexer.Pop()},
		},
		"Call": []lexer.Rule{
			lexer.Include("Whitespace"),
			lexer.Include("Atomic"),
			lexer.Include("ExprValues"),
			{Name: "CallEnd", Pattern: `\)`, Action: lexer.Pop()},
		},
		"Whitespace": []lexer.Rule{
			{Name: "Comment", Pattern: `#.*$`},
			{Name: "CommentLine", Pattern: `^[ \t]*#.*\r?\n`},
			{Name: "Whitespace", Pattern: `[ \t]+`},
		},
		"ExprValues": []lexer.Rule{
			{Name: "Integer", Pattern: INTEGER},
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
