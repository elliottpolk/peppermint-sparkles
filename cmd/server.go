// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"fmt"
	"os"
	"strings"
)

const Server = "server"

func GetConfgrAddr() string {
	if os.Getenv("CONFGR_ADDR") == "" {
		fmt.Println("CONFGR_ADDR must be set prior to usage (i.e. export CONFGR_ADDR=localhost:8080)")
		os.Exit(1)
	}

	//  ensure the address at least has http
	addr := os.Getenv("CONFGR_ADDR")
	if !strings.HasPrefix(addr, "http") {
		addr = fmt.Sprintf("http://%s", addr)
	}

	//  trim the trailing slash to allow the commands to not have to worry
	if strings.HasSuffix(addr, "/") {
		addr = strings.TrimSuffix(addr, "/")
	}

	return addr
}
