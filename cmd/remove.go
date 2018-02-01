// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package cmd

import (
	"fmt"

	"git.platform.manulife.io/go-common/log"
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

	id := context.String(flag(SecretIdFlag.Name))
	if len(id) < 1 {
		log.NewError("a valid secret ID must be provided")
		return
	}

	if _, err := del(asURL(addr, fmt.Sprintf("%s/%s", service.PathSecrets, id), "")); err != nil {
		log.Error(err, "unable to remove secrets")
		return
	}
}
