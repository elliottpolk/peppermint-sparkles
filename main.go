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
				cmd.AddrFlag,
				cmd.AppNameFlag,
				cmd.AppEnvFlag,
				cmd.DecryptFlag,
				cmd.TokenFlag,
			},
			Usage:  "retrieves all or specific secrets",
			Action: cmd.Get,
		},
		{
			Name:    "set",
			Aliases: []string{"add", "create", "new", "update"},
			Flags: []cli.Flag{
				cmd.AddrFlag,
				cmd.AppNameFlag,
				cmd.AppEnvFlag,
				cmd.SecretFlag,
				cmd.EncryptFlag,
				cmd.TokenFlag,
			},
			Usage:  "adds or updates a secret",
			Action: cmd.Set,
		},
		{
			Name:    "delete",
			Aliases: []string{"del", "rm"},
			Flags: []cli.Flag{
				cmd.AddrFlag,
				cmd.SecretIdFlag,
			},
			Usage:  "deletes the secret for the provided app name and optional environment",
			Action: cmd.Remove,
		},
		{
			Name:    "server",
			Aliases: []string{"serve"},
			Flags: []cli.Flag{
				cmd.StdListenPortFlag,
				cmd.TlsListenPortFlag,
				cmd.TlsCertFlag,
				cmd.TlsKeyFlag,
				cmd.DatastoreAddrFlag,
				cmd.DatastoreFileFlag,
				cmd.DatastoreTypeFlag,
			},
			Usage:  "start server",
			Action: cmd.Serve,
		},
	}

	app.Run(os.Args)
}
