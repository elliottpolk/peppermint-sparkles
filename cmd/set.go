// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/secret"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

func Set(context *cli.Context) error {
	addr := context.String(AddrFlag.Names()[0])
	if len(addr) < 1 {
		cli.ShowCommandHelpAndExit(context, context.Command.FullName(), 1)
		return nil
	}

	raw, f := context.String(SecretFlag.Names()[0]), context.String(SecretFileFlag.Names()[0])
	if len(raw) > 0 && len(f) > 0 {
		return cli.Exit(errors.New("only 1 input method is allowed"), 1)
	}

	//	raw should not have anything if this is true
	if len(f) > 0 {
		info, err := os.Stat(f)
		if err != nil {
			return cli.Exit(errors.Wrap(err, "uanble to access secrets file"), 1)
		}

		if info.Size() > int64(MaxData) {
			return cli.Exit(errors.New("secret must be less than 3MB"), 1)
		}

		r, err := ioutil.ReadFile(f)
		if err != nil {
			return cli.Exit(errors.Wrap(err, "unable to read in secret file"), 1)
		}

		raw = string(r)
	}

	//	if raw is still empty at this point, attempt to read in piped data
	tick := 0
	for len(raw) < 1 {
		if tick > 0 {
			return cli.Exit(errors.New("a valid secret must be specified"), 1)
		}

		r, err := pipe()
		if err != nil {
			switch err {
			case ErrNoPipe:
				return cli.Exit(errors.New("a valid secret must be specified"), 1)
			case ErrDataTooLarge:
				return cli.Exit(errors.New("secret must be less than 3MB"), 1)
			default:
				return cli.Exit(errors.Wrap(err, "unable to read piped in data"), 1)
			}
		}
		raw, tick = r, +1
	}

	s, err := secret.NewSecret(raw)
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to parse secret"), 1)
	}

	c, encrypt := &pgp.Crypter{}, context.Bool(EncryptFlag.Names()[0])
	if encrypt {
		token := context.String(TokenFlag.Names()[0])
		if len(token) < 1 {
			//	attempt to generate a token if one not provided, erroring and exiting
			//	if unable. This attempts to prevent encrypting with empty string
			t, err := crypto.NewToken()
			if err != nil {
				return cli.Exit(errors.Wrap(err, "unable to generate encryption token"), 1)
			}
			token = t
		}

		c.Token = []byte(token)

		cypher, err := c.Encrypt([]byte(s.Content))
		if err != nil {
			return cli.Exit(errors.Wrap(err, "unable to encrypt secret content"), 1)
		}

		s.Content = string(cypher)
	}

	//	convert to JSON string for sending to secrets service
	out, err := json.Marshal(s)
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to convert secret to JSON string"), 1)
	}

	res, err := send(asURL(addr, service.PathSecrets, ""), string(out))
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to send config"), 1)
	}

	//	convert from and back to JSON string to provide "prettier" formatting on print
	ugly := &secret.Secret{}
	if err := json.Unmarshal([]byte(res), &ugly); err != nil {
		log.Error(err, "unable to parse in JSON string for pretty output")
	}

	pretty, err := json.MarshalIndent(ugly, "", "   ")
	if err != nil {
		log.Error(err, "unable to marshal secret back to (prettier) JSON string")
	}

	//	ensure to display encryption token, since it may have been generated
	if encrypt {
		log.Infof("token: %s", c.Token)
	}
	log.Infof("secret:\n%s", string(pretty))

	return nil
}
