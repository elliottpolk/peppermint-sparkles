// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/crypto/pgp/pgp.go
//
package pgp

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

const PGPMessageType string = "PGP MESSAGE"

var ErrInvalidToken error = errors.New("invalid token")

type Crypter struct {
	Token []byte
}

//  Encrypt takes a key and text which attempts to encode using the OpenPGP
//  symmetrical encryption.
func (c *Crypter) Encrypt(text []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	encoder, err := armor.Encode(buf, PGPMessageType, nil)
	if err != nil {
		return nil, err
	}

	ptxtWriter, err := openpgp.SymmetricallyEncrypt(encoder, c.Token, nil, nil)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(ptxtWriter, string(text))

	ptxtWriter.Close()
	encoder.Close()

	return []byte(base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

//  Decrypt expects an OpenPGP encoded cypher, returning the decrypted results
//  of the cypher text using the provided token.
func (c *Crypter) Decrypt(cypher []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(cypher))
	if err != nil {
		return nil, err
	}

	block, err := armor.Decode(bytes.NewBuffer(decoded))
	if err != nil {
		return nil, err
	}

	readTick := 0
	details, err := openpgp.ReadMessage(block.Body, nil, func(k []openpgp.Key, s bool) ([]byte, error) {
		// 	temporary hack since this will be called several times when a given
		//	token is not valid.
		//	TODO :: review openpgp source for more info
		if readTick > 100 {
			return c.Token, ErrInvalidToken
		}

		readTick++
		return c.Token, nil
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
