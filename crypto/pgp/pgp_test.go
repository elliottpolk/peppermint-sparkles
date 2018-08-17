package pgp

import (
	"testing"
)

var (
	token1 = []byte("YTZkNDUxMjQtMDk5Ny00NTc2LTg5YTUtMzVlZGExOTlmOTk0Cg==")
	token2 = []byte("NDM1MkM4NTItQjc1QS00NzJCLUI3RDktMTBFOEZDNkMzMzRFCg==")
)

//	generated at http://fillerama.io/
const (
	filler1 = `Wow, you got that off the Internet? In my day, the Internet was 
	only used to download pornography. No, she'll probably make me do it. Oh, 
	I always feared he might run off like this. Why, why, why didn't I break his 
	legs?
    `
	filler2 = `Anyone who laughs is a communist! But I've never been to the moon! 
	You, a bobsleder!? That I'd like to see! Who are you, my warranty?! Is that a 
	cooking show? Take me to your leader! That could be 'my' beautiful soul sitting 
	naked on a couch. If I could just learn to play this stupid thing. Who said 
	that? SURE you can die! You want to die?!
	`
)

func TestEncrypt(t *testing.T) {
	c1, c2 := &Crypter{Token: token1}, &Crypter{Token: token2}
	cypher1, err := c1.Encrypt([]byte(filler1))
	if err != nil {
		t.Fatal(err)
	}

	cypher2, err := c2.Encrypt([]byte(filler1))
	if err != nil {
		t.Fatal(err)
	}

	if c1, c2 := string(cypher1), string(cypher2); c1 == c2 {
		t.Errorf("different token, same text encrypts to the same output string\n%s\n%s", c1, c2)
	}

	if c1, c2 := string(cypher1), string(cypher2); c1 == filler1 || c2 == filler2 || c1 == filler2 || c2 == filler1 {
		t.Error("text did not encrypt, resulting output is same as input")
	}
}

func TestDecrypt(t *testing.T) {
	c1, c2 := &Crypter{Token: token1}, &Crypter{Token: token2}
	cypher1, err := c1.Encrypt([]byte(filler1))
	if err != nil {
		t.Fatal(err)
	}

	cypher2, err := c1.Encrypt([]byte(filler2))
	if err != nil {
		t.Fatal(err)
	}

	//	ensure that a different token can not produce the original string
	if _, err := c2.Decrypt(cypher1); err != nil && err != ErrInvalidToken {
		t.Fatal(err)
	}

	//	tests decrypt
	ptxt1, err := c1.Decrypt(cypher1)
	if err != nil {
		t.Fatal(err)
	}

	//	validate the decryption returned the initial string
	if want, got := filler1, string(ptxt1); want != got {
		t.Errorf("\nwant %s\ngot %s\n", want, got)
	}

	ptxt2, err := c1.Decrypt(cypher2)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := filler2, string(ptxt2); want != got {
		t.Errorf("\nwant %s\ngot %s\n", want, got)
	}

	if notwant, got := filler1, string(ptxt2); notwant == got {
		t.Error("same token, different text returns the same decrypted text")
	}
}
