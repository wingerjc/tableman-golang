package main

import "flag"

func readFlags() *programOptions {
	result := &programOptions{}
	result.Web = &WebOptions{}
	flag.StringVar(&result.InputFile, "input", "", "Package file to load into main program.")
	flag.StringVar(&result.OutputFile, "output", "", "File to write results to.")
	flag.StringVar(&result.ScriptFile, "script", "", "Script file to read command from.")
	flag.BoolVar(&result.Interactive, "interact", true, "Whether to print command prompt.")
	flag.BoolVar(&result.Echo, "echo", false, "Whether to echo each commmand to output.")
	flag.StringVar(&result.CLIPrefix, "prefix", "$ ", "The prefix for command line input")

	// Web server flags
	flag.BoolVar(&result.Web.RunWeb, "web", false, "Run the program as a web server.")
	flag.StringVar(&result.Web.Addr, "web-addr", ":8080", "The local address to serve from")
	flag.StringVar(&result.Web.CertFile, "web-certfile", "", "Cert file for TLS, web-keyfile must also be defined")
	flag.StringVar(&result.Web.KeyFile, "web-keyfile", "", "Keyfile for TLS, web-certfile must also be defined")

	flag.Parse()
	return result
}

type programOptions struct {
	Web         *WebOptions
	InputFile   string
	OutputFile  string
	ScriptFile  string
	Interactive bool
	Echo        bool
	CLIPrefix   string
}

type WebOptions struct {
	RunWeb bool
	Addr   string
	// TLS settings
	CertFile string
	KeyFile  string
}
