package parser

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func parserTypeWithDefaultOptions(t interface{}) (*participle.Parser, error) {
	return participle.Build(t, participle.Lexer(fileLexer), participle.Elide("Comment", "Whitespace", "CommentLine"))
}

func TestNumberRanges(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&NumberRange{})
	assert.NoError(err)

	val := &NumberRange{}
	err = parser.ParseString("", `7`, val)
	assert.NoError(err)
	assert.Equal(7, val.Single)
	if print {
		pp.Println(val)
	}

	val2 := &NumberRange{}
	err = parser.ParseString("", `1-7`, val2)
	assert.NoError(err)
	assert.Equal(1, *val2.First)
	assert.Equal(7, val2.Last)
	if print {
		pp.Println(val2)
	}

	p2, err := parserTypeWithDefaultOptions(&RangeList{})
	assert.NoError(err)
	val3 := &RangeList{}
	err = p2.ParseString("", `1-7,4,5,9`, val3)
	assert.NoError(err)
	assert.NotNil(val3.Ranges)
	assert.Len(val3.Ranges, 4)
	if print {
		pp.Println(val3)
	}
}

func TestLabelString(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&LabelString{})
	assert.NoError(err)

	val := &LabelString{}
	expect := `"Twas brillig and the slithy toves."`
	err = parser.ParseString("", expect, val)
	assert.NoError(err)
	assert.NotNil(val.Escaped)
	assert.Equal(expect, *val.Escaped)
	if print {
		pp.Println(val)
	}

	val = &LabelString{}
	expect = `S0mething_cool-ish`
	err = parser.ParseString("", expect, val)
	assert.NoError(err)
	assert.NotNil(val.Single)
	assert.Equal(expect, *val.Single)
	if print {
		pp.Println(val)
	}
}

func TestTableRow(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&TableRow{})
	assert.NoError(err)

	// Expression value
	val := &TableRow{}
	err = parser.ParseString("", `{4}`, val)
	assert.NoError(err)
	assert.Len(val.Values, 1)
	if print {
		pp.Println(val)
	}

	// Lbael and expression value.
	val = &TableRow{}
	expect := `"qux zed"`
	err = parser.ParseString("", `zombie: {4} `+expect, val)
	assert.NoError(err)
	assert.False(val.Default)
	assert.NotNil(val.Label)
	assert.Len(val.Values, 2)
	assert.Equal(expect, *val.Values[1].StringVal)
	if print {
		pp.Println(val)
	}

	// Range selector and string value
	val = &TableRow{}
	expect = `"grok"`
	err = parser.ParseString("", `1-3,4,8: `+expect, val)
	assert.NoError(err)
	assert.False(val.Default)
	assert.NotNil(val.Numbers)
	assert.Len(val.Values, 1)
	assert.Equal(expect, *val.Values[0].StringVal)
	if print {
		pp.Println(val)
	}

	// Two string values with numeeric and label selectors.
	val = &TableRow{}
	expects := []string{`"gleebrox"`, `"molduk"`}
	err = parser.ParseString("", `6 zombie: `+expects[0]+` `+expects[1], val)
	assert.NoError(err)
	assert.False(val.Default)
	assert.NotNil(val.Label)
	assert.NotNil(val.Numbers)
	assert.Len(val.Values, 2)
	assert.Equal(expects[0], *val.Values[0].StringVal)
	assert.Equal(expects[1], *val.Values[1].StringVal)
	if print {
		pp.Println(val)
	}

	// Shortest default.
	val = &TableRow{}
	err = parser.ParseString("", `Default: "zzz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Println(val)
	}

	// Default as def and numric value.
	val = &TableRow{}
	err = parser.ParseString("", `Default 6: "zzz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Println(val)
	}

	// Default and label.
	val = &TableRow{}
	err = parser.ParseString("", `Default axes: "xyz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Println(val)
	}

	// Weighted rows
	val = &TableRow{}
	err = parser.ParseString("", `w=36: "xyz"`, val)
	assert.NoError(err)
	assert.Equal(36, val.Weight)
	if print {
		pp.Println(val)
	}

	// Count rows
	val = &TableRow{}
	err = parser.ParseString("", `c=36: "xyz"`, val)
	assert.NoError(err)
	assert.Equal(36, val.Count)
	if print {
		pp.Println(val)
	}

	// Multiple line lists.
	val = &TableRow{}
	err = parser.ParseString("", "\"xyz\" ->\n    \"asdf\" \"qwer\" ->\r\n \t\"ffff\"", val)
	assert.NoError(err)
	assert.Len(val.Values, 4)
	if print {
		pp.Println(val)
	}
}

func TestTableHeader(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&TableHeader{})
	assert.NoError(err)

	val := &TableHeader{}
	err = parser.ParseString("", `TableDef: table1`, val)
	assert.NoError(err)
	assert.Equal("table1", val.Name)
	assert.Len(val.Tags, 0)
	if print {
		pp.Println(val)
	}

	val = &TableHeader{}
	err = parser.ParseString("", "TableDef: table1 \n ~ author: franz", val)
	assert.NoError(err)
	assert.Equal("table1", val.Name)
	assert.Len(val.Tags, 1)
	if print {
		pp.Println(val)
	}

	val = &TableHeader{}
	err = parser.ParseString("", "TableDef: table1 \n ~ author: franz \n ~ \"create date\": \"10/01/2001\"", val)
	assert.NoError(err)
	assert.Equal("table1", val.Name)
	assert.Len(val.Tags, 2)
	if print {
		pp.Println(val)
	}
}

func TestTable(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Table{})
	assert.NoError(err)

	val := &Table{}
	table := `TableDef: footable
	~ license: free
	1: "asdf"
	"green"
	w=3: {4}
	foo: "once"
	Default: "red"`
	err = parser.ParseString("", table, val)
	assert.NoError(err)
	assert.Equal("footable", val.Header.Name)
	assert.Len(val.Header.Tags, 1)
	assert.Len(val.Rows, 5)
	assert.True(val.Rows[4].Default)
	if print {
		pp.Println(val)
	}
}

func TestFileHeader(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&FileHeader{})
	assert.NoError(err)

	val := &FileHeader{}
	err = parser.ParseString("", `TablePack: Main`, val)
	assert.NoError(err)
	assert.Equal("Main", val.Name.Names[0])
	if print {
		pp.Println(val)
	}

	val = &FileHeader{}
	err = parser.ParseString("", `TablePack: Foo.Bar.Baz.quz.z3d`, val)
	assert.NoError(err)
	assert.Equal("Foo", val.Name.Names[0])
	assert.Len(val.Name.Names, 5)
	assert.Equal("z3d", val.Name.Names[4])
	if print {
		pp.Println(val)
	}

	val = &FileHeader{}
	header := `TablePack: garbo
	Import: f"~/tables/1.tab"
	Import: f"c:/blah/blerf"`
	err = parser.ParseString("", header, val)
	assert.NoError(err)
	assert.Equal("garbo", val.Name.Names[0])
	assert.Len(val.Name.Names, 1)
	assert.Len(val.Imports, 2)
	if print {
		pp.Println(val)
	}

	val = &FileHeader{}
	header = `TablePack: flux-capacitor.car
	Import: f"/marty/Mcfly" As: marty
	Import: f"../doc/file.file" As: doc.Fu7ur3_space`
	err = parser.ParseString("", header, val)
	assert.NoError(err)
	assert.Equal("flux-capacitor", val.Name.Names[0])
	assert.Len(val.Name.Names, 2)
	assert.Len(val.Imports, 2)
	assert.Equal("marty", val.Imports[0].Alias.Names[0])
	assert.Len(val.Imports[1].Alias.Names, 2)
	if print {
		pp.Println(val)
	}
}

func TestTableFile(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&TableFile{})
	assert.NoError(err)

	val := &TableFile{}
	err = parser.ParseString("", `TablePack: Main`, val)
	assert.NoError(err)
	assert.Equal("Main", val.Header.Name.Names[0])
	if print {
		pp.Println(val)
	}

	val = &TableFile{}
	file := `TablePack: Main

	TableDef: colors
	"red"
	"green"
	# some note
	"blue"
	---


	TableDef: numbers
	~ size: "7"
	~ "weight max": forty-eight
	a: "One"
	5: "two"
	Default: {3}


	`
	err = parser.ParseString("", file, val)
	assert.NoError(err)
	assert.Equal("Main", val.Header.Name.Names[0])
	assert.Len(val.Tables, 2)
	if print {
		pp.Println(val)
	}

	val = &TableFile{}
	file = `TablePack: quacko
	Import: f"~/jeremy" As: something.new

	TableDef: colors
	Default w=4 1-5,8 red: "red"
	"blue"

	#----------
	TableDef: numbers
	~ size: "7"
	~ "weight max": forty-eight
	a: "One"
	5: "two"
	Default: {3}`
	err = parser.ParseString("", file, val)
	assert.NoError(err)
	assert.Equal("quacko", val.Header.Name.Names[0])
	assert.Len(val.Tables, 2)
	assert.Len(val.Tables[0].Rows, 2)
	if print {
		pp.Println(val)
	}
}

func TestRoll(t *testing.T) {
	print := false
	// t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Roll{})
	assert.NoError(err)

	val := &Roll{}
	err = parser.ParseString("", `3d8?`, val)
	assert.NoError(err)
	assert.Equal("3d8", val.RollDice)
	if print {
		pp.Println(val)
	}

	val = &Roll{}
	err = parser.ParseString("", `10d20h5.mode?`, val)
	assert.NoError(err)
	assert.Equal("10d20", val.RollDice)
	assert.Equal("h", val.RollSubset)
	assert.Equal(5, val.SubsetCount)
	assert.Equal(".mode", val.RollFuncAggr)
	assert.False(val.Print)
	if print {
		pp.Println(val)
	}

	val = &Roll{}
	err = parser.ParseString("", `10d20.+1x2.-20x1.str?`, val)
	assert.NoError(err)
	assert.Equal("10d20", val.RollDice)
	assert.Len(val.RollCountAggrs, 2)
	assert.Equal(".+", val.RollCountAggrs[0].Sign)
	assert.Equal(1, val.RollCountAggrs[0].Number)
	assert.Equal(2, val.RollCountAggrs[0].Multiplier)
	assert.True(val.Print)
	if print {
		pp.Println(val)
	}
}

func TestExpr(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Expression{})
	assert.NoError(err)

	val := &Expression{}
	err = parser.ParseString("", `{ 4}`, val)
	assert.NoError(err)
	assert.Equal(4, *val.Value.Num)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	err = parser.ParseString("", `{-8 }`, val)
	assert.NoError(err)
	assert.Equal(-8, *val.Value.Num)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	err = parser.ParseString("", `{ foo }`, val)
	assert.NoError(err)
	assert.NotNil(val.Value.Label)
	assert.Equal("foo", *val.Value.Label.Single)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	err = parser.ParseString("", `{ "hello World!" }`, val)
	assert.NoError(err)
	assert.NotNil(val.Value.Label)
	assert.Equal(`"hello World!"`, *val.Value.Label.Escaped)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	strVal := "{ add( \n\t3,\r\n	\t6) }"
	err = parser.ParseString("", strVal, val)
	assert.NoError(err)
	assert.NotNil(val.Value.Call)
	assert.False(val.Value.Call.IsTable)
	assert.Equal("add", val.Value.Call.Name.Names[0])
	assert.Len(val.Value.Call.Params, 2)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	err = parser.ParseString("", `{ !color() }`, val)
	assert.NoError(err)
	assert.NotNil(val.Value.Call)
	assert.True(val.Value.Call.IsTable)
	assert.Equal("color", val.Value.Call.Name.Names[0])
	assert.Len(val.Value.Call.Params, 0)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	err = parser.ParseString("", `{ 1d3? }`, val)
	assert.NoError(err)
	assert.NotNil(val.Value.Roll)
	assert.Equal("1d3", val.Value.Roll.RollDice)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	err = parser.ParseString("", `{ add( !primes(), 3, !test(weight, 3d12.+12x1?), oops) }`, val)
	assert.NoError(err)
	assert.NotNil(val.Value.Call)
	assert.False(val.Value.Call.IsTable)
	assert.Len(val.Value.Call.Params, 4)
	if print {
		pp.Println(val)
	}
}

func TestExprVars(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Expression{})
	assert.NoError(err)

	val := &Expression{}
	strVal := `{ @foo=6; 4}`
	err = parser.ParseString("", strVal, val)
	assert.NoError(err)
	assert.Len(val.Vars, 1)
	assert.Equal("foo", val.Vars[0].VarName.Name)
	assert.Equal(6, *val.Vars[0].AssignedValue.Num)
	assert.Equal(4, *val.Value.Num)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	strVal = `{ @one="One", @Second-Thing="two"; @Second-Thing }`
	err = parser.ParseString("", strVal, val)
	assert.NoError(err)
	assert.Len(val.Vars, 2)
	assert.Equal("\"One\"", *val.Vars[0].AssignedValue.Label.Escaped)
	assert.Equal("Second-Thing", val.Value.Variable.Name)
	if print {
		pp.Println(val)
	}

	val = &Expression{}
	strVal = "{\n\t @one=\"On\ne\",\r\n @Second-Thing=\"two\";\t\n @Second-Thing }"
	err = parser.ParseString("", strVal, val)
	assert.NoError(err)
	assert.Len(val.Vars, 2)
	assert.Equal("\"On\ne\"", *val.Vars[0].AssignedValue.Label.Escaped)
	assert.Equal("Second-Thing", val.Value.Variable.Name)
	if print {
		pp.Println(val)
	}
}

func TestGeneratorRows(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Table{})
	assert.NoError(err)

	val := &Table{}
	strVal := `TableDef: foo
	["a"]`
	err = parser.ParseString("", strVal, val)
	assert.NoError(err)
	assert.Nil(val.Rows)
	assert.Len(val.Generator.Steps, 1)
	assert.Len(val.Generator.Steps[0].Values, 1)
	if print {
		pp.Println(val)
	}

	val = &Table{}
	strVal = `TableDef: foo
	["A","2","3","4", "5", "6", "7", "8", "9", "10" , "J" ,"Q",
	"K"][" of "]["Clubs", "Spades", "Diamonds", "Hearts"]`
	err = parser.ParseString("", strVal, val)
	assert.NoError(err)
	assert.Nil(val.Rows)
	assert.Len(val.Generator.Steps, 3)
	assert.Len(val.Generator.Steps[0].Values, 13)
	assert.Len(val.Generator.Steps[1].Values, 1)
	assert.Len(val.Generator.Steps[2].Values, 4)
	assert.Equal(`"Diamonds"`, val.Generator.Steps[2].Values[2])
	if print {
		pp.Println(val)
	}
}
