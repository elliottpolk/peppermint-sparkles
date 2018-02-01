// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/secrets/main.go
//
package cmd

import (
	"bufio"
	"encoding/json"
	"io"
	"math"
	"os"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/secrets/crypto"
	"git.platform.manulife.io/oa-montreal/secrets/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/secrets/secret"
	"git.platform.manulife.io/oa-montreal/secrets/service"

	"github.com/urfave/cli"
)

func Set(context *cli.Context) {
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

	env := context.String(flag(AppEnvFlag.Name))
	if len(env) < 1 {
		log.NewError("a valid app environment must be provided")
		return
	}

	content := context.String(flag(SecretFlag.Name))
	if len(content) < 1 {
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Error(err, "unable to stat stdin")
			return
		}

		if fi.Mode()&os.ModeCharDevice != 0 || fi.Size() < 1 {
			log.NewError("a secret must be provided")
			return
		}

		buf, res := bufio.NewReader(os.Stdin), make([]byte, 0)
		for {
			in, _, err := buf.ReadLine()
			if err != nil && err == io.EOF {
				break
			}
			res = append(res, in...)

			if len(res) > (int(math.Pow10(7)) * 3) {
				log.NewError("secret data should be less than 3MB in size")
				return
			}
		}

		content = string(res)
	}

	s := &secret.Secret{
		App:     app,
		Env:     env,
		Content: content,
	}

	c := &pgp.Crypter{}
	encrypt := context.Bool(flag(EncryptFlag.Name))

	if encrypt {
		token := context.String(flag(TokenFlag.Name))
		if len(token) < 1 {
			//	attempt to generate a token if one not provided, erroring and exiting
			//	if unable. This attempts to prevent encrypting with empty string
			t, err := crypto.NewToken()
			if err != nil {
				log.Error(err, "unable to generate encryption token")
				return
			}
			token = t
		}

		c.Token = []byte(token)

		cypher, err := c.Encrypt([]byte(s.Content))
		if err != nil {
			log.Error(err, "unable to encrypt secret content")
			return
		}

		s.Content = string(cypher)
	}

	out, err := json.Marshal(s)
	if err != nil {
		log.Error(err, "unable to convert secret to JSON string")
		return
	}

	res, err := send(asURL(addr, service.PathSecrets, ""), string(out))
	if err != nil {
		log.Error(err, "unable to send config")
		return
	}

	//	convert from and back to JSON string to provide "prettier" formatting
	ugly := &secret.Secret{}
	if err := json.Unmarshal([]byte(res), &ugly); err != nil {
		log.Error(err, "unable to parse in JSON string for pretty output")
	}

	pretty, err := json.MarshalIndent(ugly, "", "   ")
	if err != nil {
		log.Error(err, "unable to marshal secret back to (prettier) JSON string")
	}

	if encrypt {
		log.Infof("token: %s", c.Token)
	}
	log.Infof("secret:\n%s", string(pretty))
}
