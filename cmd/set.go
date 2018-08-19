package cmd

import (
	"bufio"
	"io"
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

func pipe() (string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", errors.Wrap(err, "unable to stat stdin")
	}

	if fi.Mode()&os.ModeCharDevice != 0 || fi.Size() < 1 {
		return "", ErrNoPipe
	}

	buf, res := bufio.NewReader(os.Stdin), make([]byte, 0)
	for {
		in, _, err := buf.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		res = append(res, in...)

		if len(res) > MaxData {
			return "", ErrDataTooLarge
		}
	}

	return string(res), nil
}

func set(encrypt bool, token, usr, raw, addr string) (*models.Secret, error) {
	s, err := models.ParseSecret(raw)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse secret")
	}

	// ensure the secret has an ID set
	if len(s.Id) < 1 {
		s.Id = uuid.GetV4()
	}

	if encrypt {
		c := &pgp.Crypter{Token: []byte(token)}

		// encrypt the content of the secret
		cypher, err := c.Encrypt([]byte(s.Content))
		if err != nil {
			return nil, errors.Wrap(err, "unable to encrypt secret content")
		}

		// set the content to the encrypted text
		s.Content = string(cypher)
	}

	params := url.Values{
		service.UserParam: []string{usr},
		service.AppParam:  []string{s.App},
		service.EnvParam:  []string{s.Env},
		service.IdParam:   []string{s.Id},
	}

	res, err := send(asURL(addr, service.PathSecrets, params.Encode()), s.MustString())
	if err != nil {
		return nil, errors.Wrap(err, "unable to send secret")
	}

	in, err := models.ParseSecret(string(res))
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse in service response")
	}

	return in, nil
}

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

	encrypt := context.Bool(EncryptFlag.Names()[0])
	token := context.String(TokenFlag.Names()[0])
	if encrypt {
		if len(token) < 1 {
			//	attempt to generate a token if one not provided, erroring and exiting
			//	if unable. This attempts to prevent encrypting with empty string
			t, err := crypto.NewToken()
			if err != nil {
				return cli.Exit(errors.Wrap(err, "unable to generate encryption token"), 1)
			}
			token = t
		}
	}

	// get current logged in user
	u, err := user.Current()
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to retrieve current, logged-in user"), 1)
	}

	s, err := set(encrypt, token, u.Username, raw, addr)
	if err != nil {
		return cli.Exit(errors.Wrap(err, "unable to set secret"), 1)
	}

	//	ensure to display encryption token, since it may have been generated
	if encrypt {
		log.Infof(tag, "token: %s", token)
	}
	log.Infof(tag, "secret:\n%s", s.MustString())

	return nil
}
