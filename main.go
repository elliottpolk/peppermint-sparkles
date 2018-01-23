// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package main

import (
	"os"

	"git.platform.manulife.io/oa-montreal/campx/cmd"

	"github.com/urfave/cli"
)

var version string

func main() {
	app := cli.NewApp()

	app.Usage = "TODO..."
	app.Version = version

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
				cmd.SecretFlag,
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
			Name: "server",
			Aliases: []string{"serve"},
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
