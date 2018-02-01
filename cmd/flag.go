// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/cmd/flag.go
//
package cmd

import (
	"strings"

	"github.com/urfave/cli"
)

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

	EncryptFlag = cli.BoolFlag{
		Name:  "encrypt",
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
		Name:   "addr",
		Usage:  "campx service address",
		EnvVar: "CAMPX_ADDR",
	}
)

func flag(name string) string {
	s := strings.Split(name, ",")
	return strings.TrimSpace(s[0])
}
