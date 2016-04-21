// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package pgp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"

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

	ptxtWriter, err := openpgp.SymmetricallyEncrypt(encodeWriter, key, nil, nil)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(ptxtWriter, string(text))

	ptxtWriter.Close()
	encodeWriter.Close()

	return []byte(base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

//  Decrypt expects an OpenPGP encoded ciphertext, returning the decrypted results
//  of the cipher text using the provided token.
func Decrypt(token, ciphertxt []byte) ([]byte, error) {
	block, err := armor.Decode(bytes.NewBuffer(ciphertxt))
	if err != nil {
		return nil, err
	}

	details, err := openpgp.ReadMessage(block.Body, nil, func(k []openpgp.Key, s bool) ([]byte, error) {
		return token, nil
	}, nil)

	if err != nil {
		return nil, err
	}

	plaintxt, err := ioutil.ReadAll(details.UnverifiedBody)
	if err != nil {
		fmt.Println(string(plaintxt))
		return nil, err
	}

	return plaintxt, nil
}
