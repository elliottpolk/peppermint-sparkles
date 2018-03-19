// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package cmd

import (
	"encoding/json"
	"net/url"
	"path"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

func Get(context *cli.Context) error {
	addr := context.String(AddrFlag.Names()[0])
	if len(addr) < 1 {
		cli.ShowCommandHelpAndExit(context, context.Command.FullName(), 1)
		return nil
	}

	token := context.String(TokenFlag.Names()[0])
	decrypt := context.Bool(DecryptFlag.Names()[0])

	if decrypt && len(token) < 1 {
		return cli.Exit(errors.New("decrypt token must be specified in order to decrypt"), 1)
	}

	id := context.String(SecretIdFlag.Names()[0])
	app := context.String(AppNameFlag.Names()[0])
	env := context.String(AppEnvFlag.Names()[0])

	params := &url.Values{}
	from := service.PathSecrets

	if len(id) < 1 {
		if len(app) < 1 || len(env) < 1 {
			return cli.Exit(errors.New("a valid secret ID or app name / environment combo must be provided"), 1)
		} else {
			params = &url.Values{
				service.AppParam: []string{app},
				service.EnvParam: []string{env},
			}
		}
	} else {
		from = path.Join(from, id)
	}

	raw, err := retrieve(asURL(addr, from, params.Encode()))
	if err != nil && err.Error() != "no valid secret" {
		return cli.Exit(errors.Wrap(err, "unable to retrieve secret"), 1)
	}

	if len(raw) < 1 {
		return nil
	}

	//  test / validate if stored content meets the secrets model and also
	//  to allow for decryption
	s := &models.Secret{}
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return cli.Exit(errors.Wrap(err, "unable to convert string to secrets"), 1)
	}

	if decrypt {
		c := pgp.Crypter{Token: []byte(token)}
		res, err := c.Decrypt([]byte(s.Content))
		if err != nil {
			return cli.Exit(errors.Wrap(err, "unable to decrypt secret"), 1)
		}
		s.Content = string(res)
	}

	log.Infof("\n%s\n", s.MustString())

	return nil
}
