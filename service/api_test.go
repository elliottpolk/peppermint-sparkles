package service

import (
	"fmt"
	"net"
	"testing"

	"git.platform.manulife.io/go-common/log"
)

func init() {
	log.InitTester()
}

func freeport() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func TestGetId(t *testing.T) {
	want := "6f0f9805-08c6-48f2-b3c4-fe8e7c35ea4a"
	m, got, err := getId(fmt.Sprintf("secrets/%s", want))
	if err != nil {
		t.Fatal(err)
	}

	if !m {
		t.Error("expected match and getId returned false")
	}

	if got != want {
		t.Errorf("want %s\n\ngot %s\n", want, got)
	}

	noMatch, invalidId, err := getId(fmt.Sprintf("testing/get/secrets/%s-asdfqwerty", want))
	if err != nil {
		t.Fatal(err)
	}

	if noMatch {
		t.Error("expected to not match and getId returned true")
	}

	if invalidId == want {
		t.Errorf("expected no match but returned %s", invalidId)
	}
}
