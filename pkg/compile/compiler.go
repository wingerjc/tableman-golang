package compiler

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"

	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

type Compiler struct {
	parser     *parser.TableFileParser
	exprParser *parser.ExpressionParser
}

func NewCompiler() (*Compiler, error) {
	p, err := parser.GetParser()
	if err != nil {
		return nil, err
	}
	exprParser, err := parser.GetExpressionParser()
	if err != nil {
		return nil, err
	}
	return &Compiler{
		parser:     p,
		exprParser: exprParser,
	}, nil
}

func (c *Compiler) CompileFile(fileName string) (*program.Program, error) {
	code, err := c.loadFile(fileName)
	if err != nil {
		return nil, err
	}
	return c.compile(code)
}

func (c *Compiler) CompileString(code string) (*program.Program, error) {
	parsed, err := c.parser.Parse(code)
	if err != nil {
		return nil, err
	}
	key := makeKey(code)
	return c.compile(&readTable{
		key:    key,
		parsed: parsed,
	})
}
func (c *Compiler) compile(pack *readTable) (*program.Program, error) {
	tableq := make([]*readTable, 0)
	tableq = append(tableq, pack)
	tableDefs := program.NewTableMap()
	parsed := pack.parsed
	first := true
	for len(tableq) > 0 {
		// pop table
		t := tableq[0]
		tableq = tableq[1:]

		// skip processed tables
		if _, ok := tableDefs[t.key]; ok {
			continue
		}
		// Setup name <-> key conversion, qualified and non-qualified tables point to this file.
		keys := make(nameMap)
		keys[""] = t.key
		keys[parsed.Header.Name.FullName()] = t.key

		// Queue up imports
		for _, i := range t.parsed.Header.Imports {
			// filename magic to get an absolute path if we can...
			fname, err := getFileName(t.fname, i.File())
			if err != nil {
				return nil, err
			}
			// open file, get hash
			tr, err := c.loadFile(fname)
			if err != nil {
				return nil, err
			}

			// resolve table prefix
			if i.Alias != nil {
				keys[i.Alias.FullName()] = tr.key
			} else {
				keys[tr.parsed.Header.Name.FullName()] = tr.key
			}

			// don't add if already enqueued
			_, queued := tableDefs[tr.key]
			for _, x := range tableq {
				queued = queued || (x.key == tr.key)
			}
			if !queued {
				tableq = append(tableq, tr)
			}
		}

		// compile file
		pack, err := CompileTableFile(t.parsed, t.key, keys)
		if err != nil {
			return nil, err
		}
		tableDefs[t.key] = pack
		if first {
			first = false
			tableDefs[program.ROOT_PACK] = pack
		}
	}

	return program.NewProgram(tableDefs), nil
}

func CompileTableFile(parsed *parser.TableFile, key string, tableKeys nameMap) (*program.TablePack, error) {
	tables := make(map[string]*program.Table)
	for _, t := range parsed.Tables {
		compiledTable, err := CompileTable(t, &program.DefaultRandSource{}, tableKeys)
		if err != nil {
			return nil, err
		}
		tables[compiledTable.Name()] = compiledTable
	}
	return program.NewTablePack(parsed.Header.Name.FullName(), key, tables), nil
}

func (c *Compiler) parseString(code string) (*parser.TableFile, error) {
	result, err := c.parser.Parse(code)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Compiler) loadFile(fname string) (*readTable, error) {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	code := string(f)
	parsed, err := c.parser.Parse(code)
	if err != nil {
		return nil, err
	}
	return &readTable{
		fname:  fname,
		parsed: parsed,
		key:    makeKey(code),
	}, nil
}

type nameMap map[string]string

type readTable struct {
	fname  string
	parsed *parser.TableFile
	key    string
}

func makeKey(code string) string {
	hash := md5.Sum([]byte(code))
	return hex.EncodeToString(hash[:])
}

func getFileName(caller string, imported string) (string, error) {
	if filepath.IsAbs(imported) {
		return imported, nil
	}
	if len(caller) > 0 && filepath.IsAbs(caller) {
		return filepath.Join(filepath.Dir(caller), imported), nil
	}
	return filepath.Abs(imported)
}
