// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/secrets/main.go
//
package cmd

import (
	"encoding/json"
	"net/url"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/secrets/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/secrets/secret"
	"git.platform.manulife.io/oa-montreal/secrets/service"

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

	app := context.String(flag(AppNameFlag.Name))
	if len(app) < 1 {
		log.NewError("a valid app name must be provided")
		return
	}

	token := context.String(flag(TokenFlag.Name))
	decrypt := context.Bool(flag(DecryptFlag.Name))

	if decrypt && len(token) < 1 {
		log.NewError("decrypt token must be specified in order to decrypt")
		return
	}

	params := &url.Values{service.AppParam: []string{app}}
	if env := context.String(flag(AppEnvFlag.Name)); len(env) > 0 {
		params.Add(service.EnvParam, env)
	}

	raw, err := retrieve(asURL(addr, service.PathSecrets, params.Encode()))
	if err != nil && err.Error() != "no valid secret" {
		log.Error(err, "unable to retrieve secret")
		return
	}

	if len(raw) < 1 {
		return
	}

	//  test / validate if stored content meets the secrets model and also
	//  to allow for decryption
	secrets := make([]*secret.Secret, 0)
	if err := json.Unmarshal([]byte(raw), &secrets); err != nil {
		log.Error(err, "unable to convert string to secrets")
		return
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
}
