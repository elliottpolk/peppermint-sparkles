// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("usage: %s <command> [args]\n\n", os.Args[0])

		fmt.Println("Available commands:")
		fmt.Printf("\t%s\t\tstarts confgr server\n", Server)
		fmt.Printf("\t%s\t\tretrieves the available app configs\n", List)
		fmt.Printf("\t%s\t\tretrieves the available config\n", Get)
		fmt.Printf("\t%s\t\tadds a new config\n", Set)
		fmt.Printf("\t%s\t\tdeletes the specified config\n", Remove)

		if len(os.Args) > 1 {
			switch os.Args[1] {
			case Get:
				fmt.Printf("Arguments for %s:\n", Get)
				fmt.Printf("\t%s\tapp name to retrieve respective config\n", AppFlag)
				fmt.Printf("\t%s\tdecrypt the config\n", DecryptFlag)
				fmt.Printf("\t%s\ttoken to decrypt config with\n", TokenFlag)

			case Set:
				fmt.Printf("Arguments for %s:\n", Set)
				fmt.Printf("\t%s\tapp name to be set\n", AppFlag)
				fmt.Printf("\t%s\tconfig to be written\n", CfgFlag)
				fmt.Printf("\t%s\tencrypt config\n", EncryptFlag)

			case Remove:
				fmt.Printf("Arguments for %s:\n", Remove)
				fmt.Printf("\t%s\tapp name to be removed\n", AppFlag)
			}
		}

		fmt.Println()
		os.Exit(0)
	}
}
