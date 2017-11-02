// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"encoding/json"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/log"
	"github.com/elliottpolk/confgr/service"

	"github.com/urfave/cli"
)

func Set(context *cli.Context) {
	context.Command.VisibleFlags()

	addr := context.String(Simplify(AddrFlag.Name))
	if len(addr) < 1 {
		if err := cli.ShowCommandHelp(context, context.Command.FullName()); err != nil {
			log.Error(err, "unable to display help")
		}
		return
	}

	app := context.String(Simplify(AppFlag.Name))
	if len(app) < 1 {
		if err := cli.ShowCommandHelp(context, context.Command.FullName()); err != nil {
			log.Error(err, "unable to display help")
		}
		return
	}

	cfg := &config.Config{
		App:         app,
		Environment: context.String(Simplify(EnvFlag.Name)),
		Value:       context.String(Simplify(ConfFlag.Name)),
	}

	token := context.String(Simplify(TokenFlag.Name))
	encrypt := context.Bool(Simplify(EncryptFlag.Name))
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
