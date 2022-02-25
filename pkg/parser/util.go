package parser

import (
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/k0kubun/pp"
)

var typeString map[lexer.TokenType]string

func setupTypeTable() {
	if len(typeString) == 0 {
		typeString = make(map[lexer.TokenType]string)
		for k, v := range fileLexer.Symbols() {
			typeString[v] = k
		}
	}
}

// PrintTokens can be used for debug to print out up to
// the first `count` tokens parsed from the input stream.
func PrintTokens(in string, count int) {
	setupTypeTable()
	l, err := fileLexer.LexString("", in)
	pp.Println(err)
	t, _ := l.Next()
	for i := 0; int(t.Type) != -1 && i < count; i++ {
		pp.Println(typeString[t.Type] + " " + t.Value)
		t, _ = l.Next()
	}
}
