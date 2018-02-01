// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/secrets/cmd/flag.go
//
package cmd

import "gopkg.in/urfave/cli.v2"

var (
	AppNameFlag = cli.StringFlag{
		Name:  "a, app, app-name",
		Usage: "app name of secret",
	}

	AppEnvFlag = cli.StringFlag{
		Name:  "e, env, app-env",
		Usage: "environment of secret (e.g. PROD, DEV, TEST, etc.)",
	}

	SecretIdFlag = cli.StringFlag{
		Name:  "id, secret-id",
		Usage: "generated ID of secret",
	}

	SecretFlag = cli.StringFlag{
		Name:  "s, secret",
		Usage: "secret to be stored",
	}

	SecretFileFlag = cli.StringFlag{
		Name:  "f, secret-file",
		Usage: "filepath to secret",
	}

	EncryptFlag = cli.BoolFlag{
		Name:  "encrypt",
		Value: true,
		Usage: "encrypt secrets",
	}

	DecryptFlag = cli.BoolFlag{
		Name:  "decrypt",
		Usage: "decrypt secrets",
	}

	TokenFlag = cli.StringFlag{
		Name:  "t, token",
		Usage: "token used to encrypt / decrypt secrets",
	}

	AddrFlag = cli.StringFlag{
		Name:    "addr",
		Usage:   "secrets service address",
		EnvVars: []string{"SECRETS_ADDR"},
	}
)
