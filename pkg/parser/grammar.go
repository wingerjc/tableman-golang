package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	NONE_EXPR_T  = ValueExprType(0)
	ROLL_EXPR_T  = ValueExprType(1)
	LABEL_EXPR_T = ValueExprType(2)
	NUM_EXPR_T   = ValueExprType(3)
	TABLE_EXPR_T = ValueExprType(4)
	FUNC_EXPR_T  = ValueExprType(5)
	VAR_EXPR_T   = ValueExprType(6)
)

var (
	DEFAULT_LEXER = participle.Lexer(fileLexer)
	DEFAULT_ELIDE = participle.Elide("Comment", "Whitespace", "CommentLine")
	STR_EXPR_TYPE = map[ValueExprType]string{
		NONE_EXPR_T:  "None",
		ROLL_EXPR_T:  "Roll",
		LABEL_EXPR_T: "Label",
		NUM_EXPR_T:   "Number",
		TABLE_EXPR_T: "Table",
		FUNC_EXPR_T:  "Function",
		VAR_EXPR_T:   "Variable",
	}
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

func (i *ImportStatement) File() string {
	return i.FileName[2 : len(i.FileName)-1]
}

type Table struct {
	Pos       lexer.Position
	Header    *TableHeader       `parser:"@@"`
	Rows      []*TableRow        `parser:"((EOL @@)+"`
	Generator *GeneratorTableRow `parser:"| (EOL@@))"`
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

type GeneratorTableRow struct {
	Steps []*GeneratorStep `parser:"@@ (EOL? @@)*"`
}

type GeneratorStep struct {
	Values []string `parser:"GenStart @String (ListDelimiter EOL? @String)* GenEnd"`
}

func (s *GeneratorStep) StrVal(index int) string {
	v := s.Values[index]
	return v[1 : len(v)-1]
}

type TableRow struct {
	Pos     lexer.Position
	Default bool         `parser:"(@Default?"`
	Weight  int          `parser:"(WeightMarker @Number)?"`
	Count   int          `parser:"(CountMarker @Number)?"`
	Numbers *RangeList   `parser:"@@?"`
	Label   *LabelString `parser:"@@? ':')?"`
	Values  []*RowItem   `parser:"@@+"`
}

type RowItem struct {
	Pos        lexer.Position
	StringVal  *string     `parser:"(@String"`
	Expression *Expression `parser:"| @@)(ExtendLine EOL)?"`
}

func (r *RowItem) String() string {
	if r.StringVal == nil {
		return ""
	}
	l := len(*r.StringVal)
	return (*r.StringVal)[1 : l-1]
}

type LabelString struct {
	Pos     lexer.Position
	Single  *string `parser:"@TableName"`
	Escaped *string `parser:"| @String"`
}

func (l *LabelString) String() string {
	if l.Single == nil {
		sLen := len(*l.Escaped)
		return (*l.Escaped)[1 : sLen-1]
	}
	return *l.Single
}

func (l *LabelString) IsLabel() bool {
	return l.Single == nil
}

type RangeList struct {
	Pos    lexer.Position
	Ranges []*NumberRange `parser:"@@ (ListDelimiter @@)*"`
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
	RollDice       string           `parser:"@Roll"`
	RollSubset     string           `parser:"(@RollSubset"`
	SubsetCount    int              `parser:"@Number)? ("`
	RollFuncAggr   string           `parser:"@RollFuncAggr"`
	RollCountAggrs []*RollCountAggr `parser:"|(@@+))?"`
	Print          bool             `parser:"@RollCast? RollEnd"`
}

type RollCountAggr struct {
	Sign       string `parser:"@RollCountSign"`
	Number     int    `parser:"@Number"`
	Multiplier int    `parser:"RollCountMultiplier @Number"`
}

type Expression struct {
	Pos   lexer.Position
	Vars  []*VariableDef `parser:"ExprStart EOL? (@@ (ListDelimiter EOL? @@)* EndVarList EOL?)?"`
	Value *ValueExpr     `parser:"@@ EOL? ExprEnd"`
}

type VariableDef struct {
	VarName       *VarName   `parser:"@@ VarAssign"`
	AssignedValue *ValueExpr `parser:"@@"`
}

type Call struct {
	IsTable bool              `parser:"@TableCallSignal?"`
	Name    ExtendedTableName `parser:"@@ CallStart EOL?"`
	Params  []*ValueExpr      `parser:"@@? (ListDelimiter EOL? @@)* EOL? CallEnd"`
}

type ValueExpr struct {
	Roll     *Roll        `parser:"@@"`
	Num      *int         `parser:"| (@Number | @Integer)"`
	Call     *Call        `parser:"| @@"`
	Label    *LabelString `parser:"| @@"`
	Variable *VarName     `parser:"| @@"`
	exprType ValueExprType
}

type ValueExprType int

func (v *ValueExpr) GetType() ValueExprType {
	if v.exprType != NONE_EXPR_T {
		return v.exprType
	} else if v.Roll != nil {
		v.exprType = ROLL_EXPR_T
	} else if v.Num != nil {
		v.exprType = NUM_EXPR_T
	} else if v.Label != nil {
		v.exprType = LABEL_EXPR_T
	} else if v.Call != nil {
		if v.Call.IsTable {
			v.exprType = TABLE_EXPR_T
		} else {
			v.exprType = FUNC_EXPR_T
		}
	} else if v.Variable != nil {
		v.exprType = VAR_EXPR_T
	}
	return v.exprType
}

func (v *ValueExpr) GetStringType() string {
	return STR_EXPR_TYPE[v.GetType()]
}

type VarName struct {
	Name string `parser:"VarPrefix @TableName"`
}
type ExtendedTableName struct {
	Names []string `parser:" @TableName (PkgDelimiter @TableName)*"`
}

func (n *ExtendedTableName) PackageName() string {
	return strings.Join(n.Names[:(len(n.Names)-1)], ".")
}

func (n *ExtendedTableName) TableName() string {
	return n.Names[len(n.Names)-1]
}

func (n *ExtendedTableName) FullName() string {
	return strings.Join(n.Names, ".")
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
			{Name: "WeightMarker", Pattern: `w=`},
			{Name: "CountMarker", Pattern: `c=`},
			{Name: "ExtendLine", Pattern: `->`},
			{Name: "TableBarrier", Pattern: `--(-+)`},
			{Name: "FilePath", Pattern: `f\"(([A-Za-z]:)|~|(\.\.?))?(/|(\\)+).*\"`},
			lexer.Include("Atomic"),
			{Name: "TableDelimiter", Pattern: `:`},
			{Name: "RangeDash", Pattern: `-`},
			{Name: "TagStart", Pattern: `~`},
			{Name: "GenStart", Pattern: `\[`},
			{Name: "GenEnd", Pattern: `]`},
		},
		"Atomic": []lexer.Rule{
			{Name: "TableName", Pattern: IDENTIFIER},
			{Name: "Roll", Pattern: NATURAL_NUMBER + `d` + NATURAL_NUMBER, Action: lexer.Push("Roll")},
			{Name: "CallStart", Pattern: `\(`, Action: lexer.Push("Call")},
			{Name: "ExprStart", Pattern: `{`, Action: lexer.Push("Expr")},
			{Name: "String", Pattern: `"(\\"|[^"])*"`},
			lexer.Include("NumberRule"),
			{Name: "TableCallSignal", Pattern: `\!`},
			{Name: "VarPrefix", Pattern: `@`},
			{Name: "PkgDelimiter", Pattern: `\.`},
			{Name: "ListDelimiter", Pattern: `,`},
			{Name: "EOL", Pattern: `\r?\n`},
		},
		"Roll": []lexer.Rule{
			{Name: "RollSubset", Pattern: `(l|h)`},
			{Name: "RollFuncAggr", Pattern: `\.(min|max|sum|avg|mode)`},
			{Name: "RollCountSign", Pattern: `\.[+-]`},
			{Name: "RollCountMultiplier", Pattern: `x`},
			{Name: "RollCast", Pattern: `\.str`},
			{Name: "RollEnd", Pattern: `\?`, Action: lexer.Pop()},
			lexer.Include("NumberRule"),
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
		"NumberRule": []lexer.Rule{
			{Name: "Number", Pattern: NATURAL_NUMBER},
		},
	})
)
