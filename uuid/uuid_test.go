package uuid

import (
	"regexp"
	"testing"
)

//  since the UUID is fairly random, the test checks for the correct format in
//  order to "validate" the returned token
func TestGetV4(t *testing.T) {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")

	if token := GetV4(); !r.MatchString(token) {
		t.Errorf("invalid UUID produced: got %s\n", token)
	}
}
