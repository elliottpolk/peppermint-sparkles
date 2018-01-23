// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package pgp

import "testing"

const token = "YTZkNDUxMjQtMDk5Ny00NTc2LTg5YTUtMzVlZGExOTlmOTk0Cg=="

//	generated at http://fillerama.io/
const filler = `Wow, you got that off the Internet? In my day, the Internet was 
	only used to download pornography. No, she'll probably make me do it. Oh, 
	I always feared he might run off like this. Why, why, why didn't I break his 
	legs?
    `

func TestDecrypt(t *testing.T) {
	//	tests encrypt
	ciphertxt, err := Encrypt([]byte(token), []byte(filler))
	if err != nil {
		t.Fatalf("unable to encrypt filler text: %v\n", err)
	}

	//	tests decrypt
	plaintxt, err := Decrypt([]byte(token), ciphertxt)
	if err != nil {
		t.Errorf("unable to decrypt text: %v\n", err)
		t.FailNow()
	}

	if string(plaintxt) != filler {
		t.Errorf("decrypt failed\n")
		t.Errorf("expected: %s\n", filler)
		t.Errorf("got:      %s\n", string(plaintxt))
		t.FailNow()
	}
}
