package crypto

import (
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Crypter interface {
	Encrypt(text []byte) ([]byte, error)
	Decrypt(cypher []byte) ([]byte, error)
}

func NewToken() (string, error) {
	token := uuid.New().String()
	if len(token) < 1 {
		return "", errors.New("UUID produced empty string")
	}

	return base64.StdEncoding.EncodeToString([]byte(token)), nil
}
