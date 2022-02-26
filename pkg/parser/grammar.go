package parser

import (
	"strconv"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	// NoneExprT the type value for an untyped value expression.
	NoneExprT ValueExprType = 0

	// RollExprT the type value for an roll value expression.
	RollExprT ValueExprType = 1

	// LabelExprT the type value for a label/string value expression.
	LabelExprT ValueExprType = 2

	// NumExprT the type value for a numeric value expression.
	NumExprT ValueExprType = 3

	// TableExprT the type value for a table call value expression.
	TableExprT ValueExprType = 4

	// FuncExprT the type value for a function call value expression.
	FuncExprT ValueExprType = 5

	// VarExprT the type value for a variable value expression.
	VarExprT ValueExprType = 6
)

var (
	// DefaultLexer is a default lexer for the tableman language.
	DefaultLexer = participle.Lexer(fileLexer)
	// DefaultElide the default list of tokens to elide from AST parsing.
	DefaultElide = participle.Elide("Comment", "Whitespace", "CommentLine")
	// ExprTypeStr a human readable map of expression value types for debugging.
	ExprTypeStr = map[ValueExprType]string{
		NoneExprT:  "None",
		RollExprT:  "Roll",
		LabelExprT: "Label",
		NumExprT:   "Number",
		TableExprT: "Table",
		FuncExprT:  "Function",
		VarExprT:   "Variable",
	}
)

// TableFile is an AST node that incorporates a whole table source file.
//
//  Pattern:
//    <EOL>*
//    <FileHeader>
//    (<EOL>+ <Table> <TableBarrier>?)*
// `TableBarrier` is at least 3 dashes on its own line.
type TableFile struct {
	Pos    lexer.Position
	Header *FileHeader `parser:"EOL* @@"`
	Tables []*Table    `parser:"(EOL+ @@? TableBarrier?)*"`
}

// FileHeader is an AST node that describes a table file.
//
//  Pattern:
//    TablePack: <ExtendedTableName>
//    (<EOL> <ImportStatement>)*
//
//  Example:
//    TablePack: foo.bar.baz
//    Import f"~/quz/util.tman" As: q.util
type FileHeader struct {
	Pos     lexer.Position
	Name    *ExtendedTableName `parser:"PkgStart @@"`
	Imports []*ImportStatement `parser:"(EOL @@)*"`
}

// ImportStatement is an AST node that denotes an imported file.
//
//  Pattern:
//    Import: <FilePath> (As: <ExtendedTableName>)?
//
//  Example:
//    Import f"~/quz/util.tman" As: q.util
//    Import f"~/quz/cards.tman"
type ImportStatement struct {
	Pos      lexer.Position
	FileName string             `parser:"Import @FilePath"`
	Alias    *ExtendedTableName `parser:"(PackAlias @@)?"`
}

// File returns the actual file name and not the parsed token.
func (i *ImportStatement) File() string {
	return i.FileName[2 : len(i.FileName)-1]
}

// Table is an AST node that denotes a single table.
//
// It can be provided either a list of table rows or a single generator row to
// programatically create rows from.
//
//  Pattern:
//    <TableHeader>
//    (
//        (<EOL> <TableRow>)+
//      | (<EOL> <GeneratorTableRow>)
//    )
type Table struct {
	Pos       lexer.Position
	Header    *TableHeader       `parser:"@@"`
	Rows      []*TableRow        `parser:"((EOL @@)+"`
	Generator *GeneratorTableRow `parser:"| (EOL@@))"`
}

// TableHeader is an AST node that denotes meta information about a table.
//
// Tags are currently transferred in compilation, but not otherwise accessible in the
// execution engine. This may change in the future.
//
//  Pattern:
//    TableDef: <TableName>
//    (<EOL>+ <Tag>)*
//
//  Example:
//    TableDef: GreekNames
//    ~ something: something-else
//    ~ "With spaces": "needs quotes"
type TableHeader struct {
	Pos  lexer.Position
	Name string `parser:"TableStart @TableName"`
	Tags []*Tag `parser:"(EOL+ @@)*"`
}

// Tag is an AST node that denotes a meta tag.
//
// Useful for tagging author, source, copyright/license, or other information.
//
//  Pattern:
//    ~ <Label>: <Label>
//
//  Example:
//    ~ foo: bar
type Tag struct {
	Pos   lexer.Position
	Key   LabelString `parser:"TagStart @@ TableDelimiter"`
	Value LabelString `parser:"@@"`
}

// GeneratorTableRow is an AST node that denotes an ordered list of row generation steps.
//
//  Pattern:
//    <GeneratorStep> (<EOL>? <GeneratorStep>)*
//
//  Example:
//    ["1", "2", "3"]["a", "b", "c"]
type GeneratorTableRow struct {
	Steps []*GeneratorStep `parser:"@@ (EOL? @@)*"`
}

// GeneratorStep is an AST node that denotes a list of generation targets.
//
//  Pattern:
//    [ <Label> (, <EOL>? <Label>)* ]
//
//  Example:
//    ["x", "y", "z"]
type GeneratorStep struct {
	Values []string `parser:"GenStart @String (ListDelimiter EOL? @String)* GenEnd"`
}

// StrVal returns the actual string to be generated at the given index. Convenience method.
func (s *GeneratorStep) StrVal(index int) string {
	v := s.Values[index]
	return v[1 : len(v)-1]
}

// TableRow is an AST node that denotes a single table row.
//
// Each row can be the default row for label and index lookups.
// It can also have optional weighted lookup values, counts for deck draws
// and be indepentently assigned ranges for index lookups and a label.
//
// A row needs at least one value, but all values will be concatenated as strings.
//
//  Pattern:
//    Default? (w=<Number>)? (c=<number>)? <RangeList>? <Label>? :?
//    (<RowItem> (-> <EOL>)? )+
//
//  Example:
//    w=5 Hard-TH: "th" ->
//        "'"
type TableRow struct {
	Pos     lexer.Position
	Default bool         `parser:"(@Default?"`
	Weight  int          `parser:"(WeightMarker @Number)?"`
	Count   int          `parser:"(CountMarker @Number)?"`
	Numbers *RangeList   `parser:"@@?"`
	Label   *LabelString `parser:"@@? ':')?"`
	Values  []*RowItem   `parser:"(@@ (ExtendLine EOL)? )+"`
}

// RowItem is an AST node that denotes a single value to be concatenated in a row.
//
// The line extension `->` can be used to shorten longer lines for readability.
//
//  Pattern:
//    (<Label> | <Expression>)
type RowItem struct {
	Pos        lexer.Position
	StringVal  *string     `parser:"(@String"`
	Expression *Expression `parser:"| @@)"`
}

// String returns the wrapped passed string. Convenience method.
func (r *RowItem) String() string {
	if r.StringVal == nil {
		return ""
	}
	l := len(*r.StringVal)
	return (*r.StringVal)[1 : l-1]
}

// LabelString is an AST node that can be either a string or a name.
//
// Except for in table rows, if your text follows the TableName format
// you can omit double quotes for simplicity and clarity.
//
// TableName is of the pattern: [a-zA-Z][a-zA-Z0-9\-_]*
//
// Strings can include carriage returns and span lines.
//
//  Pattern:
//    <TableName> | "<string>"
type LabelString struct {
	Pos     lexer.Position
	Single  *string `parser:"@TableName"`
	Escaped *string `parser:"| @String"`
}

// String returns the string or label value. Convenience method.
func (l *LabelString) String() string {
	if l.Single == nil {
		sLen := len(*l.Escaped)
		return (*l.Escaped)[1 : sLen-1]
	}
	return *l.Single
}

// IsLabel returns whether this value can be processed as a label.
func (l *LabelString) IsLabel() bool {
	return l.Single == nil
}

// RangeList is an AST node that denotes a list of integers or number ranges
//
//  Pattern:
//    <NumberRange> (, <NumberRange>)*
//
//  Example:
//    1, 2, 5-9, 8
type RangeList struct {
	Pos    lexer.Position
	Ranges []*NumberRange `parser:"@@ (ListDelimiter @@)*"`
}

// NumberRange is an AST node that denotes a single number or a range of numbers.
//
//  Pattern:
//    <Number>
//  | <Number>-<Number>
//
//  Example:
//    9-16
type NumberRange struct {
	Pos    lexer.Position
	First  *int `parser:"((@Number'-'"`
	Last   int  `parser:"@Number)"`
	Single int  `parser:"| @Number)"`
}

// Roll is an AST node that denotes a dice roll expression.
//
//  Pattern:
//    <Number>d<Number>
//    ( (l | h) <Number>)?
//    (
//        . (min | max | sum | avg | mode | median)
//      | (<RollCountAggr>)+
//    )?
//    .str? <?>
//
//  Example:
//    9d5.h6.median
//
//  <?> is a literal quesiton mark
type Roll struct {
	Pos            lexer.Position
	RollDice       string           `parser:"@Roll"`
	RollSubset     string           `parser:"(@RollSubset"`
	SubsetCount    int              `parser:"@Number)? ("`
	RollFuncAggr   string           `parser:"@RollFuncAggr"`
	RollCountAggrs []*RollCountAggr `parser:"|(@@+))?"`
	Print          bool             `parser:"@RollCast? RollEnd"`
}

// Dice returns the count of dice and how many sides are on each die.
func (r *Roll) Dice() (count int, sides int, err error) {
	nums := strings.Split(r.RollDice, "d")
	count, err = strconv.Atoi(nums[0])
	if err != nil {
		return
	}
	sides, err = strconv.Atoi(nums[1])
	return
}

// FnAggr is a convenience method for extracting a normalized function aggr name.
// If there is no function aggr, an empty string is returned.
func (r *Roll) FnAggr() string {
	if len(r.RollFuncAggr) == 0 {
		return ""
	}
	return r.RollFuncAggr[1:]
}

// RollCountAggr is an AST node that multiplies th
//
//  Pattern:
//    . (+ | -) <Number> (x <Number>)?
//
//  Example:
//    .+20x3
type RollCountAggr struct {
	Sign       string `parser:"@RollCountSign"`
	Number     int    `parser:"@Number"`
	Multiplier int    `parser:"(RollCountMultiplier @Number)?"`
}

// FinalMult returns the final signed multiplier for the node
func (r *RollCountAggr) FinalMult() int {
	if r.Multiplier == 0 {
		r.Multiplier = 1
	}
	if r.Sign[1:] == "-" {
		return -1 * r.Multiplier
	}
	return r.Multiplier
}

// Expression is an AST node for an expression that can have variables defined within it.
//
//  Pattern:
//    {
//       <EOL>?
//       (<VariableDef> (, <EOL>? <VariableDef>)* ; <EOL>? )?
//       <ValueExpr> <EOL>?
//    }
//
//  Example:
//    { @foo=8, @bar=1d8?; add(@foo, @bar) }
type Expression struct {
	Pos   lexer.Position
	Vars  []*VariableDef `parser:"ExprStart EOL? (@@ (ListDelimiter EOL? @@)* EndVarList EOL?)?"`
	Value *ValueExpr     `parser:"@@ EOL? ExprEnd"`
}

// VariableDef is an AST node for defining a variable.
//
//  Pattern:
//    <VarName> = <ValueExpr>
//
//  Example:
//    @foo=9
type VariableDef struct {
	VarName       *VarName   `parser:"@@ VarAssign"`
	AssignedValue *ValueExpr `parser:"@@"`
}

// Call is an AST node for calling table or function.
//
// Table calls are delineated by starting with an exclamation point.
//
//  Pattern:
//    !? <ExtendedTableName> <(> (<ValueExpr> (, <EOL>? <ValueExpr>)* )? <)>
//
//  Example:
//    !CardDeck(deck, shuffle)
type Call struct {
	IsTable bool              `parser:"@TableCallSignal?"`
	Name    ExtendedTableName `parser:"@@ CallStart EOL?"`
	Params  []*ValueExpr      `parser:"(@@ (ListDelimiter EOL? @@)* )? EOL? CallEnd"`
}

// ValueExpr is an AST node for expressions that can return values.
// Reused in a few places as a general building block.
//
//  Pattern:
//    <Roll>
//  | <Number>
//  | <Call>
//  | <LabelString>
//  | <VarName>
type ValueExpr struct {
	Roll     *Roll        `parser:"@@"`
	Num      *int         `parser:"| (@Number | @Integer)"`
	Call     *Call        `parser:"| @@"`
	Label    *LabelString `parser:"| @@"`
	Variable *VarName     `parser:"| @@"`
	exprType ValueExprType
}

// ValueExprType is an enum type for denoting a ValueExpr stored type.
type ValueExprType int

// GetType resolves and caches the type of this ValueExpr.
func (v *ValueExpr) GetType() ValueExprType {
	if v.exprType != NoneExprT {
		return v.exprType
	} else if v.Roll != nil {
		v.exprType = RollExprT
	} else if v.Num != nil {
		v.exprType = NumExprT
	} else if v.Label != nil {
		v.exprType = LabelExprT
	} else if v.Call != nil {
		if v.Call.IsTable {
			v.exprType = TableExprT
		} else {
			v.exprType = FuncExprT
		}
	} else if v.Variable != nil {
		v.exprType = VarExprT
	}
	return v.exprType
}

// GetStringType returns the string version of the type value for debugging.
func (v *ValueExpr) GetStringType() string {
	return ExprTypeStr[v.GetType()]
}

// VarName is an AST node for a variable name.
//
//  Pattern:
//    @ <TableName>
type VarName struct {
	Name string `parser:"VarPrefix @TableName"`
}

// ExtendedTableName is an AST node for a package definition or
// call to a package prefixed table.
//
//  Pattern:
//    <TableName> (. <TableName>)*
//
//  Example:
//    foo.bar.baz
type ExtendedTableName struct {
	Names []string `parser:" @TableName (PkgDelimiter @TableName)*"`
}

// PackageName returns all but the last name from the list of names.
func (n *ExtendedTableName) PackageName() string {
	return strings.Join(n.Names[:(len(n.Names)-1)], ".")
}

// TableName returns the last name from the list of names.
func (n *ExtendedTableName) TableName() string {
	return n.Names[len(n.Names)-1]
}

// FullName returns a single string value for the whole list.
// Useful for package definitions.
func (n *ExtendedTableName) FullName() string {
	return strings.Join(n.Names, ".")
}

const (
	naturalNumberPat = `([1-9][0-9]*)`
	wholeNumberPat   = `(0|([1-9][0-9]*))`
	integerPat       = `(0|(-?[1-9][0-9]*))`
	identifierPat    = `[a-zA-Z][a-zA-Z0-9\-_]*`
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
			{Name: "TableName", Pattern: identifierPat},
			{Name: "Roll", Pattern: naturalNumberPat + `d` + naturalNumberPat, Action: lexer.Push("Roll")},
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
			{Name: "RollFuncAggr", Pattern: `\.(min|max|sum|avg|mode|median)`},
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
			{Name: "Integer", Pattern: integerPat},
		},
		"NumberRule": []lexer.Rule{
			{Name: "Number", Pattern: naturalNumberPat},
		},
	})
)
