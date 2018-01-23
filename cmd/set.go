// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package cmd

import (
	"encoding/json"

	"git.platform.manulife.io/oa-montreal/campx/config"
	"git.platform.manulife.io/oa-montreal/campx/log"
	"git.platform.manulife.io/oa-montreal/campx/service"

	"github.com/urfave/cli"
)

func Set(context *cli.Context) {
	context.Command.VisibleFlags()

	addr := context.String(flag(AddrFlag.Name))
	if len(addr) < 1 {
		if err := cli.ShowCommandHelp(context, context.Command.FullName()); err != nil {
			log.Error(err, "unable to display help")
		}
		return
	}

	app := context.String(flag(AppFlag.Name))
	if len(app) < 1 {
		if err := cli.ShowCommandHelp(context, context.Command.FullName()); err != nil {
			log.Error(err, "unable to display help")
		}
		return
	}

	cfg := &config.Config{
		App:         app,
		Environment: context.String(flag(EnvFlag.Name)),
		Value:       context.String(flag(SecretFlag.Name)),
	}

	token := context.String(flag(TokenFlag.Name))
	encrypt := context.Bool(flag(EncryptFlag.Name))
	if encrypt {
		t, err := cfg.Encrypt(token)
		if err != nil {
			log.Error(err, "unable to encrypt config")
			return
		}

		token = t
	}

	out, err := json.Marshal(cfg)
	if err != nil {
		log.Error(err, "unable to marshal config")
		return
	}

	res, err := send(asURL(addr, service.PathSet, ""), string(out))
	if err != nil {
		log.Error(err, "unable to send config")
		return
	}

	if encrypt {
		log.Infof("token: %s", token)
	}

	log.Infof("config:\n%s", res)
}
