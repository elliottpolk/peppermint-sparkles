// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"git.platform.manulife.io/oa-montreal/campx/log"
	"git.platform.manulife.io/oa-montreal/campx/service"

	"github.com/urfave/cli"
)

func Remove(context *cli.Context) {
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

	params := &url.Values{service.AppParam: []string{app}}

	env := context.String(flag(EnvFlag.Name))
	if len(env) > 0 {
		params.Add(service.EnvParam, env)
	} else {
		var res string

		fmt.Printf("remove all configs for %s? ", app)
		fmt.Scanf("%s", &res)

		res = strings.ToLower(res)

		if len(res) > 0 && res[0:1] != "y" {
			return
		}
	}

	if _, err := retrieve(asURL(addr, service.PathRemove, params.Encode())); err != nil {
		log.Error(err, "unable to remove config")
		return
	}
}
