// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"strings"

	"github.com/urfave/cli"
)

var (
	AppFlag = cli.StringFlag{
		Name:  "a, app",
		Usage: "app name of config",
	}

	EnvFlag = cli.StringFlag{
		Name:  "e, env",
		Usage: "environment of configuration (e.g. PROD, DEV, TEST, etc.)",
	}

	ConfFlag = cli.StringFlag{
		Name:  "c, cfg, conf, config",
		Usage: "config to be stored",
	}

	EncryptFlag = cli.BoolFlag{
		Name:  "encrypt",
		Usage: "encrypt configuration",
	}

	DecryptFlag = cli.BoolFlag{
		Name:  "decrypt",
		Usage: "decrypt configuration",
	}

	TokenFlag = cli.StringFlag{
		Name:  "t, token",
		Usage: "token used to encrypt / decrypt configuration",
	}

	AddrFlag = cli.StringFlag{
		Name:   "addr",
		Usage:  "confgr service address",
		EnvVar: "CONFGR_ADDR",
	}

	TlsCertFlag = cli.StringFlag{
		Name:   "tls-cert",
		Usage:  "TLS certificate file",
		EnvVar: "CONFGR_TLS_CERT",
	}

	TlsKeyFlag = cli.StringFlag{
		Name:   "tls-key",
		Usage:  "TLS key file",
		EnvVar: "CONFGR_TLS_KEY",
	}

	DatastoreFlag = cli.StringFlag{
		Name:   "dsf, datastore-file",
		Value:  "/var/lib/confgr/confgr.db",
		Usage:  "name / location of file for storing content",
		EnvVar: "CONFGR_DS_FILE",
	}
)

func Simplify(name string) string {
	s := strings.Split(name, ",")
	return strings.TrimSpace(s[len(s)-1])
}
