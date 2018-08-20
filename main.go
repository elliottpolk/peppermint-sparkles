package main

import (
	"os"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/cmd"

	"gopkg.in/urfave/cli.v2"
)

var version string

func main() {
	log.Init(version)

	app := cli.App{
		Copyright: "Copyright © 2018 Manulife",
		Usage:     "Server and client for managing super special secrets 🦄",
		Version:   version,
		Commands: []*cli.Command{
			cmd.Get,
			cmd.Set,
			cmd.Remove,
			cmd.Serve,
		},
	}

	app.Run(os.Args)
}
