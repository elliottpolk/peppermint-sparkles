// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package cmd

import (
	"encoding/json"
	"net/url"

	"git.platform.manulife.io/oa-montreal/campx/config"
	"git.platform.manulife.io/oa-montreal/campx/log"
	"git.platform.manulife.io/oa-montreal/campx/service"

	"github.com/urfave/cli"
)

func Get(context *cli.Context) {
	context.Command.VisibleFlags()

	addr := context.String(flag(AddrFlag.Name))
	if len(addr) < 1 {
		if err := cli.ShowCommandHelp(context, context.Command.FullName()); err != nil {
			log.Error(err, "unable to display help")
		}
		return
	}

	token := context.String(flag(TokenFlag.Name))
	decrypt := context.Bool(flag(DecryptFlag.Name))

	if decrypt && len(token) < 1 {
		log.NewError("decrypt token must be specified in order to decrypt")
		return
	}

	params := &url.Values{}
	if env := context.String(flag(EnvFlag.Name)); len(env) > 0 {
		params.Add(service.EnvParam, env)
	}

	if app := context.String(flag(AppFlag.Name)); len(app) > 0 {
		params.Add(service.AppParam, app)
	}

	raw, err := retrieve(asURL(addr, service.PathFind, params.Encode()))
	if err != nil {
		log.Error(err, "unable to retrieve config")
		return
	}

	cfgs := make([]*config.Config, 0)
	if err := json.Unmarshal([]byte(raw), &cfgs); err != nil {
		log.Error(err, "unable to convert string to configs")
		return
	}

	for _, cfg := range cfgs {
		if decrypt && len(cfg.Value) > 0 {
			if err := cfg.Decrypt(token); err != nil {
				log.Error(err, "unable to decrypt config")
			}
		}

		log.Infof("\n%s\n", cfg.MustString())
	}

}
