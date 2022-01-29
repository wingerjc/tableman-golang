package parser

import (
	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func parserTypeWithDefaultOptions(t interface{}) (*participle.Parser, error) {
	return participle.Build(t, participle.Lexer(fileLexer), participle.Elide("Comment", "Whitespace"))
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
		pp.Print(val)
	}

	val2 := &NumberRange{}
	err = parser.ParseString("", `1-7`, val2)
	assert.NoError(err)
	assert.Equal(1, *val2.First)
	assert.Equal(7, val2.Last)
	if print {
		pp.Print(val2)
	}

	p2, err := parserTypeWithDefaultOptions(&RangeList{})
	assert.NoError(err)
	val3 := &RangeList{}
	err = p2.ParseString("", `1-7,4,5,9`, val3)
	assert.NoError(err)
	assert.NotNil(val3.First)
	assert.Len(val3.Rest, 3)
	if print {
		pp.Print(val3)
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
		pp.Print(val)
	}

	val = &RowLabel{}
	expect = `S0mething_cool-ish`
	err = parser.ParseString("", expect, val)
	assert.NoError(err)
	assert.NotNil(val.Single)
	assert.Equal(expect, *val.Single)
	if print {
		pp.Print(val)
	}
}

func TestTableRow(t *testing.T) {
	print := true
	t.Parallel()
	assert := assert.New(t)
	parser, err := parserTypeWithDefaultOptions(&TableRow{})
	assert.NoError(err)

	val := &TableRow{}
	err = parser.ParseString("", `{4}`, val)
	assert.NoError(err)
	assert.Len(val.Values, 1)
	if print {
		pp.Print(val)
	}

	val = &TableRow{}
	expect := `"qux zed"`
	err = parser.ParseString("", `zombie: {4} `+expect, val)
	assert.NoError(err)
	assert.False(val.Default)
	assert.NotNil(val.Label)
	assert.Len(val.Values, 2)
	assert.Equal(expect, *val.Values[1].StringVal)
	if print {
		pp.Print(val)
	}

	val = &TableRow{}
	expect = `"grok"`
	err = parser.ParseString("", `1-3,4,8: `+expect, val)
	assert.NoError(err)
	assert.False(val.Default)
	assert.NotNil(val.Numbers)
	assert.Len(val.Values, 1)
	assert.Equal(expect, *val.Values[0].StringVal)
	if print {
		pp.Print(val)
	}

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
		pp.Print(val)
	}

	val = &TableRow{}
	err = parser.ParseString("", `d: "zzz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Print(val)
	}

	val = &TableRow{}
	err = parser.ParseString("", `def 6: "zzz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Print(val)
	}

	val = &TableRow{}
	err = parser.ParseString("", `default axes: "xyz"`, val)
	assert.NoError(err)
	assert.True(val.Default)
	if print {
		pp.Print(val)
	}

	assert.Fail("foo")
}
