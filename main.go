// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package main

import (
	"os"

	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/cmd"

	"gopkg.in/urfave/cli.v2"
)

var version string

func main() {
	app := cli.App{
		Copyright: "Copyright © 2018 Manulife",
		Usage:     "Server and client for managing super special secrets 🦄",
		Version:   version,
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "get",
				Aliases: []string{"ls", "list"},
				Flags: []cli.Flag{
					&cmd.AddrFlag,
					&cmd.AppNameFlag,
					&cmd.AppEnvFlag,
					&cmd.DecryptFlag,
					&cmd.TokenFlag,
				},
				Usage:  "retrieves secrets",
				Action: cmd.Get,
			},
			&cli.Command{
				Name:    "set",
				Aliases: []string{"add", "create", "new", "update"},
				Flags: []cli.Flag{
					&cmd.AddrFlag,
					&cmd.SecretFlag,
					&cmd.SecretFileFlag,
					&cmd.EncryptFlag,
					&cmd.TokenFlag,
				},
				Usage:  "adds or updates a secret",
				Action: cmd.Set,
			},
			&cli.Command{
				Name:    "delete",
				Aliases: []string{"del", "rm"},
				Flags: []cli.Flag{
					&cmd.AddrFlag,
					&cmd.SecretIdFlag,
				},
				Usage:  "deletes a secret",
				Action: cmd.Remove,
			},
			&cli.Command{
				Name:    "server",
				Aliases: []string{"serve"},
				Flags: []cli.Flag{
					&cmd.StdListenPortFlag,
					&cmd.TlsListenPortFlag,
					&cmd.TlsCertFlag,
					&cmd.TlsKeyFlag,
					&cmd.DatastoreAddrFlag,
					&cmd.DatastoreFileFlag,
					&cmd.DatastoreTypeFlag,
				},
				Usage:  "start server",
				Action: cmd.Serve,
			},
		},
	}

	app.Run(os.Args)
}
