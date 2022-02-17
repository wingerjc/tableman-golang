package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	compiler "github.com/wingerjc/tableman-golang/pkg/compile"
	"github.com/wingerjc/tableman-golang/pkg/program"
)

const (
	CLIPrefix = "$"
)

func main() {
	opt := readFlags()
	app, err := NewApp(opt)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	err = app.cliLoop()
	app.out.Flush()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}

type App struct {
	out         *bufio.Writer
	cmd         *bufio.Scanner
	prog        *program.Program
	compiler    *compiler.Compiler
	interactive bool
	echo        bool
}

func NewApp(opt *programOptions) (*App, error) {
	app := &App{}
	app.cmd = bufio.NewScanner(os.Stdin)
	app.interactive = true && opt.Interactive
	if len(opt.ScriptFile) > 0 {
		app.interactive = false
		f, err := os.Open(opt.ScriptFile)
		if err != nil {
			return nil, fmt.Errorf("could not open script file '%s': %w", opt.ScriptFile, err)
		}
		app.cmd = bufio.NewScanner(f)
	}
	app.out = bufio.NewWriter(os.Stdout)
	if len(opt.OutputFile) > 0 {
		f, err := os.Create(opt.OutputFile)
		if err != nil {
			return nil, fmt.Errorf("could not open output file '%s': %w", opt.OutputFile, err)
		}
		app.out = bufio.NewWriter(f)
	}
	c, err := compiler.NewCompiler()
	if err != nil {
		return nil, fmt.Errorf("could not create compiler: %w", err)
	}
	app.compiler = c
	if len(opt.InputFile) > 0 {
		if err := app.loadProgram(opt.InputFile); err != nil {
			return nil, err
		}
	}
	app.echo = opt.Echo
	return app, nil
}

func (app *App) cliLoop() error {
	if app.interactive {
		app.PF("%s ", CLIPrefix)
	}
	for app.cmd.Scan() {
		command := strings.TrimLeft(app.cmd.Text(), " \t")
		if app.echo {
			app.PF("--> %s\n", command)
		}
		first := strings.ToLower(strings.Split(command, " ")[0])
		rest := command[len(first):]
		switch first {
		// exit commands
		case "bye":
			fallthrough
		case "exit":
			fallthrough
		case "q":
			fallthrough
		case "quit":
			if app.interactive {
				app.P("Goodbye!")
			}
			return nil
		// exec commands
		case "e":
			fallthrough
		case "exec":
			if err := app.executeStatement(rest); err != nil {
				if !app.interactive {
					return err
				}
				app.P("Error executing statement: %s\n", err.Error())
			}
		// Load new program
		case "l":
			fallthrough
		case "load":
			if err := app.loadProgram(rest); err != nil {
				if !app.interactive {
					return err
				}
				app.P("Could not load program '%s': %s\n", rest, err.Error())
			}
		default:
			if app.interactive {
				app.P("Could not understand command '%s'\n", first)
			}
		}
		app.Flush()
		if app.interactive {
			app.PF("%s ", CLIPrefix)
		}
	}
	return nil
}

func (app *App) loadProgram(fname string) error {
	newProg, err := app.compiler.CompileFile(fname)
	if err != nil {
		return err
	}
	app.prog = newProg
	return nil
}

func (app *App) executeStatement(code string) error {
	if app.prog == nil {
		return fmt.Errorf("could not execute statement: no program loaded")
	}
	comp, err := app.compiler.CompileExpression(code)
	if err != nil {
		return err
	}
	result, err := app.prog.Eval(comp)
	if err != nil {
		return err
	}
	if result.MatchType(program.STRING_RESULT) {
		app.P("%s\n", result.StringVal())
	} else {
		app.P("%d\n", result.IntVal())
	}
	return nil
}

func (app *App) P(format string, vals ...interface{}) error {
	_, err := fmt.Fprintf(app.out, format, vals...)
	return err
}

func (app *App) Flush() error {
	return app.out.Flush()
}

func (app *App) PF(format string, vals ...interface{}) error {
	if err := app.P(format, vals...); err != nil {
		return err
	}
	return app.Flush()
}
