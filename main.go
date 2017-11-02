// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"os"

	"github.com/elliottpolk/confgr/cmd"

	"github.com/urfave/cli"
)

var Version string

func main() {
	app := cli.NewApp()

	app.Usage = ""
	app.Version = Version

	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"ls", "list"},
			Flags: []cli.Flag{
				cmd.AppFlag,
				cmd.EnvFlag,
				cmd.DecryptFlag,
				cmd.TokenFlag,
				cmd.AddrFlag,
			},
			Usage:  "retrieves all or specific configs",
			Action: cmd.Get,
		},
		{
			Name:    "set",
			Aliases: []string{"add"},
			Flags: []cli.Flag{
				cmd.AppFlag,
				cmd.EnvFlag,
				cmd.ConfFlag,
				cmd.EncryptFlag,
				cmd.TokenFlag,
				cmd.AddrFlag,
			},
			Usage:  "adds or updates a config",
			Action: cmd.Set,
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "rm"},
			Flags: []cli.Flag{
				cmd.AppFlag,
				cmd.EnvFlag,
				cmd.AddrFlag,
			},
			Usage:  "deletes the config for the provided app name and optional environment",
			Action: cmd.Remove,
		},
		{
			Name: "serve",
			Flags: []cli.Flag{
				cmd.TlsCertFlag,
				cmd.TlsKeyFlag,
				cmd.DatastoreFlag,
			},
			Usage:  "starts confgr server",
			Action: cmd.Serve,
		},
	}

	app.Run(os.Args)
}
