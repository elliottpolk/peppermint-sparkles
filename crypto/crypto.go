package crypto

import (
	"encoding/base64"

	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/uuid"

	"github.com/pkg/errors"
)

type Crypter interface {
	Encrypt(text []byte) ([]byte, error)
	Decrypt(cypher []byte) ([]byte, error)
}

func NewToken() (string, error) {
	token := uuid.GetV4()
	if len(token) < 1 {
		return "", errors.New("UUID produced empty string")
	}

	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}
