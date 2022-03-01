package main

import (
	"log"
)

func main() {
	opt := readFlags()

	// Run as a web server
	if opt.Web.RunWeb {
		cfg := NewServerConfig()

		s, err := NewServer(cfg)
		if err != nil {
			log.Fatal(err)
		}

		err = s.Run()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	// Run as a CLI REPL app,
	app, err := NewApp(opt)
	if err != nil {
		log.Fatal(err)
	}
	err = app.CLILoop()
	app.out.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
