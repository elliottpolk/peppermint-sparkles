package models

import (
	"testing"
)

func TestSecret(t *testing.T) {
	const sample string = `{
 "id": "6f0f9805-08c6-48f2-b3c4-fe8e7c35ea4a",
 "app_name": "dummy",
 "env": "test",
 "content": "notSuperS3cret"
}`

	// tests parsing of secret
	s, err := ParseSecret(sample)
	if err != nil {
		t.Fatal(err)
	}

	// tests stringer of secret
	if want, got := sample, s.MustString(); want != got {
		t.Errorf("want: %s\n\ngot: %s", want, got)
	}
}
