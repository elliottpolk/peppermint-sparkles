package cmd

import "gopkg.in/urfave/cli.v2"

var (
	AppNameFlag = cli.StringFlag{
		Name:    "app-name",
		Aliases: []string{"a", "app"},
		Usage:   "app name of secret",
	}

	AppEnvFlag = cli.StringFlag{
		Name:    "app-env",
		Aliases: []string{"e", "env"},
		Usage:   "environment of secret (e.g. PROD, DEV, TEST, etc.)",
	}

	SecretIdFlag = cli.StringFlag{
		Name:    "secret-id",
		Aliases: []string{"id", "sid"},
		Usage:   "generated ID of secret",
	}

	SecretFlag = cli.StringFlag{
		Name:    "secret",
		Aliases: []string{"s"},
		Usage:   "secret to be stored",
	}

	SecretFileFlag = cli.StringFlag{
		Name:    "secret-file",
		Aliases: []string{"f"},
		Usage:   "filepath to secret",
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
		Name:    "token",
		Aliases: []string{"t", "tok"},
		Usage:   "token used to encrypt / decrypt secrets",
	}

	AddrFlag = cli.StringFlag{
		Name:    "addr",
		Usage:   "secrets service address",
		EnvVars: []string{"PSPARKLES_ADDR"},
	}
)
