package compiler

import (
	"github.com/alecthomas/participle/v2"
	"github.com/wingerjc/tableman-golang/pkg/parser"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

type Compiler struct {
	parser     *participle.Parser
	exprParser *participle.Parser
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

func (c *Compiler) CompileFile(fileName string) (*program.Program, []error) {
	return nil, nil
}

func (c *Compiler) CompileString(code string) (*program.Program, []error) {
	errs := make([]error, 0, 1)
	ast, err := c.parseString(code)
	if err != nil {
		errs = append(errs, err)
	}
	prog := program.NewProgram()
	pack, compErrs := program.TablePackFromAST(ast)
	errs = append(errs, compErrs...)
	files := make([]string, 0)
	processed := make(map[string]bool, 0)
	for _, i := range pack.Imports {
		files = append(files, i.FileName)
	}

	// TODO: Need file <-> pack relationship for alias functionality.
	for len(files) > 0 {
		f := files[0]
		files = files[1:]

		if in, ok := processed[f]; ok && in {
			continue
		}
		// load file into mem,
		// compileString
		pack = &program.TablePack{}
		// add new files to file queue
		processed[f] = true
		if err := prog.AddPack(pack); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}
	return prog, nil
}

func (c *Compiler) compileString(code string) (*program.TablePack, []error) {
	return nil, nil
}

func (c *Compiler) parseString(code string) (*parser.TableFile, error) {
	result := &parser.TableFile{}
	err := c.parser.ParseString("", code, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
