package pgp

import (
	"testing"
)

func TestDecrypt(t *testing.T) {
	const (
		token    string = "MTMwNmYyZTEtMmI0Zi00NzcyLTgyMDMtODA0NGQ1N2U4ZmI0Cg=="
		cipher   string = "LS0tLS1CRUdJTiBQR1AgTUVTU0FHRS0tLS0tCgp3eDRFQndNSXErKzBHU0NJUU5KZ3VvM2JJVGxDN3BYYVFpWGpzamhjSzF2UzRBSGtaVzBpTUdnZHUwQVdBTFphClAwVGJQT0Z3S2VBRTRNcmhseTdnTHVLZko3UkI0Sm5rMDdvTVdHejFKeExveDZTR3ZNTUpHZUJkNHpOdGI0SzAKNUxZcDRBbmlSWUtjQU9CRzRiUzM0TFRnWitDcDVMVXRYVVIyOHNxcFNQTExDb01MZE5QaWpkQ0g1ZUUvTmdBPQo9SEpidgotLS0tLUVORCBQR1AgTUVTU0FHRS0tLS0t"
		expected string = "some not so random testing text"
	)

	got, err := Decrypt([]byte(token), []byte(cipher))
	if err != nil {
		t.Errorf("unable to encrypt data: %v\n", err)
	}

	if string(got) != expected {
		t.Errorf("expected %s - got %s\n", expected, got)
	}
}
