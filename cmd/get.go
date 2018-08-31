package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var (
	Get = &cli.Command{
		Name:    "get",
		Aliases: []string{"ls", "list"},
		Flags: []cli.Flag{
			&AddrFlag,
			&AppNameFlag,
			&AppEnvFlag,
			&SecretIdFlag,
			&DecryptFlag,
			&TokenFlag,
			&InsecureFlag,
		},
		Usage: "retrieves secrets",
		Action: func(context *cli.Context) error {
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

			params := &url.Values{
				service.AppParam: []string{context.String(AppNameFlag.Names()[0])},
				service.EnvParam: []string{context.String(AppEnvFlag.Names()[0])},
			}

			insecure := context.Bool(InsecureFlag.Names()[0])

			s, err := get(decrypt, insecure, token, addr, context.String(SecretIdFlag.Names()[0]), params)
			if err != nil {
				return cli.Exit(errors.Wrap(err, "unable to retrieve secert"), 1)
			}

			log.Infof("\n%s\n", s.MustString())
			return nil
		},
	}
)

func get(decrypt, insecure bool, token, addr, id string, params *url.Values) (*models.Secret, error) {
	if len(id) < 1 {
		return nil, errors.New("a valid secret ID must be provided")
	}

	if len(params.Get(service.AppParam)) < 1 {
		return nil, errors.New("a valid secret app name must be provided")
	}

	if len(params.Get(service.EnvParam)) < 1 {
		return nil, errors.New("a valid secret environment must be provided")
	}

	raw, err := retrieve(asURL(addr, fmt.Sprintf("%s/%s", service.PathSecrets, id), params.Encode()), insecure)
	if err != nil {
		if err.Error() == "no valid secret" {
			return nil, err
		}
		return nil, errors.Wrap(err, "unable to retrieve secret")
	}

	if len(raw) < 1 {
		return nil, errors.New("no valid secret")
	}

	//  test / validate if stored content meets the secrets model and also
	//  to allow for decryption
	s := &models.Secret{}
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return nil, errors.Wrap(err, "unable to convert string to secrets")
	}

	if decrypt {
		c := pgp.Crypter{Token: []byte(token)}
		res, err := c.Decrypt([]byte(s.Content))
		if err != nil {
			return nil, errors.Wrap(err, "unable to decrypt secret")
		}
		s.Content = string(res)
	}

	return s, nil
}
