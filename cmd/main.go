package main

import (
	"os"
)

func main() {
	opt := readFlags()
	app, err := NewApp(opt)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	err = app.CLILoop()
	app.out.Flush()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
