// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package cmd

import (
	"fmt"

	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

func Remove(context *cli.Context) error {
	addr := context.String(AddrFlag.Names()[0])
	if len(addr) < 1 {
		cli.ShowCommandHelpAndExit(context, context.Command.FullName(), 1)
		return nil
	}

	id := context.String(SecretIdFlag.Names()[0])
	if len(id) < 1 {
		return cli.Exit(errors.New("a valid secret ID must be provided"), 1)
	}

	if _, err := del(asURL(addr, fmt.Sprintf("%s/%s", service.PathSecrets, id), "")); err != nil {
		return cli.Exit(errors.Wrap(err, "unable to remove secrets"), 1)
	}

	return nil
}
