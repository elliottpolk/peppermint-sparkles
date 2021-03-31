package main

import (
	"os"

	"gopkg.in/urfave/cli.v2"
)

var version string

func main() {
	app := cli.App{
		Copyright: "Copyright © 2018",
		Usage:     "Server and client for managing super special secrets 🦄",
		Version:   version,
		Commands: []*cli.Command{
			Get,
			Set,
			Remove,
			Serve,
		},
	}

	app.Run(os.Args)
}
