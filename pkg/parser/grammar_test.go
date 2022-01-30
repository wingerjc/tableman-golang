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
	t.Parallel()
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
	assert.NotNil(val3.First)
	assert.Len(val3.Rest, 3)
	if print {
		pp.Println(val3)
	}
}

func TestRowLabel(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&RowLabel{})
	assert.NoError(err)

	val := &RowLabel{}
	expect := `"Twas brillig and the slithy toves."`
	err = parser.ParseString("", expect, val)
	assert.NoError(err)
	assert.NotNil(val.Escaped)
	assert.Equal(expect, *val.Escaped)
	if print {
		pp.Println(val)
	}

	val = &RowLabel{}
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
	t.Parallel()
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
	err = parser.ParseString("", `d: "zzz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Println(val)
	}

	// Default as def and numric value.
	val = &TableRow{}
	err = parser.ParseString("", `def 6: "zzz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Println(val)
	}

	// Default and label.
	val = &TableRow{}
	err = parser.ParseString("", `default axes: "xyz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Println(val)
	}

	// Weighted rows
	val = &TableRow{}
	err = parser.ParseString("", `w=36: "xyz"`, val)
	assert.NoError(err)
	assert.Equal("w=36", val.Weight)
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
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&TableHeader{})
	assert.NoError(err)

	val := &TableHeader{}
	err = parser.ParseString("", `t: table1`, val)
	assert.NoError(err)
	assert.Equal("table1", val.Name)
	assert.Len(val.Tags, 0)
	if print {
		pp.Println(val)
	}

	val = &TableHeader{}
	err = parser.ParseString("", "T: table1 \n ~ author: franz", val)
	assert.NoError(err)
	assert.Equal("table1", val.Name)
	assert.Len(val.Tags, 1)
	if print {
		pp.Println(val)
	}

	val = &TableHeader{}
	err = parser.ParseString("", "T: table1 \n ~ author: franz \n ~ \"create date\": \"10/01/2001\"", val)
	assert.NoError(err)
	assert.Equal("table1", val.Name)
	assert.Len(val.Tags, 2)
	if print {
		pp.Println(val)
	}
}

func TestTable(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Table{})
	assert.NoError(err)

	val := &Table{}
	table := `t: footable
	~ license: free
	1: "asdf"
	"green"
	w=3: {4}
	foo: "once"
	d: "red"`
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
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&FileHeader{})
	assert.NoError(err)

	val := &FileHeader{}
	err = parser.ParseString("", `TablePack Main`, val)
	assert.NoError(err)
	assert.Equal("Main", val.RootPackageName)
	if print {
		pp.Println(val)
	}

	val = &FileHeader{}
	err = parser.ParseString("", `TablePack Foo.Bar.Baz.quz.z3d`, val)
	assert.NoError(err)
	assert.Equal("Foo", val.RootPackageName)
	assert.Len(val.SubPackages, 4)
	assert.Equal("z3d", val.SubPackages[3])
	if print {
		pp.Println(val)
	}

	val = &FileHeader{}
	header := `TablePack garbo
	import f"~/tables/1.tab"
	IMPORT f"c:/blah/blerf"`
	err = parser.ParseString("", header, val)
	assert.NoError(err)
	assert.Equal("garbo", val.RootPackageName)
	assert.Len(val.SubPackages, 0)
	assert.Len(val.Imports, 2)
	if print {
		pp.Println(val)
	}
}

func TestTableFile(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&TableFile{})
	assert.NoError(err)

	val := &TableFile{}
	err = parser.ParseString("", `TablePack Main`, val)
	assert.NoError(err)
	assert.Equal("Main", val.Header.RootPackageName)
	if print {
		pp.Println(val)
	}

	val = &TableFile{}
	file := `TablePack Main

	t: colors
	"red"
	"green"
	# some note
	"blue"
	---


	t: numbers
	~ size: "7"
	~ "weight max": forty-eight
	a: "One"
	5: "two"
	d: {3}


	`
	err = parser.ParseString("", file, val)
	assert.NoError(err)
	assert.Equal("Main", val.Header.RootPackageName)
	assert.Len(val.Tables, 2)
	if print {
		pp.Println(val)
	}

	val = &TableFile{}
	file = `TablePack quacko
	import f"~/jeremy"

	Table: colors
	d w=4 1-5,8 red: "red"
	"blue"

	#----------
	table: numbers
	~ size: "7"
	~ "weight max": forty-eight
	a: "One"
	5: "two"
	d: {3}`
	err = parser.ParseString("", file, val)
	assert.NoError(err)
	assert.Equal("quacko", val.Header.RootPackageName)
	assert.Len(val.Tables, 2)
	assert.Len(val.Tables[0].Rows, 2)
	if print {
		pp.Println(val)
	}
}

func TestRoll(t *testing.T) {
	print := false
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&Roll{})
	assert.NoError(err)

	val := &Roll{}
	err = parser.ParseString("", `3d8`, val)
	assert.NoError(err)
	assert.Equal("3d8", val.RollDice)
	if print {
		pp.Println(val)
	}

	val = &Roll{}
	err = parser.ParseString("", `10d20h5.mode`, val)
	assert.NoError(err)
	assert.Equal("10d20", val.RollDice)
	assert.Equal("h5", val.RollSubset)
	assert.Equal(".mode", val.RollFuncAggr)
	if print {
		pp.Println(val)
	}

	val = &Roll{}
	err = parser.ParseString("", `10d20.+1x2.-20x1`, val)
	assert.NoError(err)
	assert.Equal("10d20", val.RollDice)
	assert.Len(val.RollCountAggrs, 2)
	assert.Equal(".+1x2", val.RollCountAggrs[0])
	if print {
		pp.Println(val)
	}
}
