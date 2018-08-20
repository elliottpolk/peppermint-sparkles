package cmd

import (
	"net"

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
