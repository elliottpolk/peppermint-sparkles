// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package pgp

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

const PGPMessageType string = "PGP MESSAGE"

//  Encrypt takes a key and text which attempts to encode using the OpenPGP
//  symmetrical encryption.
func Encrypt(key, text []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	encodeWriter, err := armor.Encode(buf, PGPMessageType, nil)
	if err != nil {
		return nil, err
	}
	defer encodeWriter.Close()

	ptxtWriter, err := openpgp.SymmetricallyEncrypt(encodeWriter, key, nil, nil)
	if err != nil {
		return nil, err
	}
	defer ptxtWriter.Close()

	fmt.Fprintf(ptxtWriter, string(text))

	return []byte(base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}
