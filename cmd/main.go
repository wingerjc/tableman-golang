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
	err = app.cliLoop()
	app.out.Flush()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
}
