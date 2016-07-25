// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/elliottpolk/confgr/cmd"
	"github.com/elliottpolk/confgr/server"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		flag.Usage()
	}

	for _, a := range args[1:] {
		if a == "-h" || a == "-help" || a == "--help" {
			flag.Usage()
		}
	}

	if args[1] != cmd.Server {
		switch args[1] {
		case cmd.List:
			if err := cmd.ListCfgs(); err != nil {
				fmt.Printf("unable to retrieve app listings: %v\n", err)
				os.Exit(1)
			}

		case cmd.Get:
			if err := cmd.GetCfg(args); err != nil {
				fmt.Printf("unable to retrieve app config: %v\n", err)
				os.Exit(1)
			}

		case cmd.Set:
			if err := cmd.SetCfg(args); err != nil {
				fmt.Printf("unable to set app config: %v\n", err)
				os.Exit(1)
			}

		case cmd.Remove:
			if err := cmd.RemoveCfg(args); err != nil {
				fmt.Printf("unable to remove app config: %v\n", err)
				os.Exit(1)
			}

		default:
			flag.Usage()
		}

		return
	}

	server.Start()
}
