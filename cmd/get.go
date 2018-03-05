// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package cmd

import (
	"encoding/json"
	"net/url"

	"gitlab.manulife.com/go-common/log"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/crypto/pgp"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/models"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

func Get(context *cli.Context) error {
	addr := context.String(AddrFlag.Names()[0])
	if len(addr) < 1 {
		cli.ShowCommandHelpAndExit(context, context.Command.FullName(), 1)
		return nil
	}

	app := context.String(AppNameFlag.Names()[0])
	if len(app) < 1 {
		return cli.Exit(errors.New("a valid app name must be provided"), 1)
	}

	token := context.String(TokenFlag.Names()[0])
	decrypt := context.Bool(DecryptFlag.Names()[0])

	if decrypt && len(token) < 1 {
		return cli.Exit(errors.New("decrypt token must be specified in order to decrypt"), 1)
	}

	params := &url.Values{service.AppParam: []string{app}}
	if env := context.String(AppEnvFlag.Names()[0]); len(env) > 0 {
		params.Add(service.EnvParam, env)
	}

	raw, err := retrieve(asURL(addr, service.PathSecrets, params.Encode()))
	if err != nil && err.Error() != "no valid secret" {
		return cli.Exit(errors.Wrap(err, "unable to retrieve secret"), 1)
	}

	if len(raw) < 1 {
		return nil
	}

	//  test / validate if stored content meets the secrets model and also
	//  to allow for decryption
	secrets := make([]*models.Secret, 0)
	if err := json.Unmarshal([]byte(raw), &secrets); err != nil {
		return cli.Exit(errors.Wrap(err, "unable to convert string to secrets"), 1)
	}

	for _, s := range secrets {
		if decrypt {
			c := pgp.Crypter{Token: []byte(token)}
			res, err := c.Decrypt([]byte(s.Content))
			if err != nil {
				log.Error(err, "unable to decrypt secret")
				continue
			}
			s.Content = string(res)
		}

		log.Infof("\n%s\n", s.MustString())
	}
	return nil
}
