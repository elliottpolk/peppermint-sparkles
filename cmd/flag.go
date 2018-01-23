// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package cmd

import (
	"strings"

	"github.com/urfave/cli"
)

var (
	AppFlag = cli.StringFlag{
		Name:  "a, app",
		Usage: "app name of secrets",
	}

	EnvFlag = cli.StringFlag{
		Name:  "e, env",
		Value: "default",
		Usage: "environment of secret (e.g. PROD, DEV, TEST, etc.)",
	}

	SecretFlag = cli.StringFlag {
		Name: "s, secret",
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

	TlsCertFlag = cli.StringFlag{
		Name:   "tls-cert",
		Usage:  "TLS certificate file for HTTPS",
		EnvVar: "CAMPX_TLS_CERT",
	}

	TlsKeyFlag = cli.StringFlag{
		Name:   "tls-key",
		Usage:  "TLS key file for HTTPS",
		EnvVar: "CAMPX_TLS_KEY",
	}

	DatastoreFlag = cli.StringFlag{
		Name:   "dsf, datastore-file",
		Value:  "/var/lib/confgr/campx.db",
		Usage:  "name / location of file for storing secrets",
		EnvVar: "CAMPX_DS_FILE",
	}
)

func flag(name string) string {
	s := strings.Split(name, ",")
	return strings.TrimSpace(s[0])
}
