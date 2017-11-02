// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/elliottpolk/confgr/log"
	"github.com/elliottpolk/confgr/service"

	"github.com/urfave/cli"
)

func Remove(context *cli.Context) {
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

	params := &url.Values{service.AppParam: []string{app}}

	env := context.String(Simplify(EnvFlag.Name))
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
