package main

import "flag"

func readFlags() *programOptions {
	result := &programOptions{}
	flag.StringVar(&result.InputFile, "input", "", "Package file to load into main program.")
	flag.StringVar(&result.OutputFile, "output", "", "File to write results to.")
	flag.StringVar(&result.ScriptFile, "script", "", "Script file to read command from.")
	flag.BoolVar(&result.Interactive, "interact", true, "Whether to print command prompt.")
	flag.BoolVar(&result.Echo, "echo", false, "Whether to echo each commmand to output.")

	flag.Parse()
	return result
}

type programOptions struct {
	InputFile   string
	OutputFile  string
	ScriptFile  string
	Interactive bool
	Echo        bool
}
