package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/uuid"

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

	s, err := models.ParseSecret(raw)
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

	//	generate and set a uuid for uniqueness
	id := uuid.GetV4()

	//	convert to JSON string for sending to secrets service
	out, err := json.Marshal(s)
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to convert secret to JSON string"), 1)
	}

	u, err := user.Current()
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to retrieve current, logged-in user"), 1)
	}

	params := url.Values{
		service.UserParam: []string{u.Username},
		service.AppParam:  []string{s.App},
		service.EnvParam:  []string{s.Env},
		service.IdParam:   []string{id},
	}

	res, err := send(asURL(addr, service.PathSecrets, params.Encode()), string(out))
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to send config"), 1)
	}

	//	convert from and back to JSON string to provide "prettier" formatting on print
	ugly := &models.Secret{}
	if err := json.Unmarshal([]byte(res), &ugly); err != nil {
		log.Error(tag, err, "unable to parse in JSON string for pretty output")
	}

	pretty, err := json.MarshalIndent(ugly, "", "   ")
	if err != nil {
		log.Error(tag, err, "unable to marshal secret back to (prettier) JSON string")
	}

	//	ensure to display encryption token, since it may have been generated
	if encrypt {
		log.Infof(tag, "token: %s", c.Token)
	}
	log.Infof(tag, "uuid:  %s", id)
	log.Infof(tag, "secret:\n%s", string(pretty))

	return nil
}
