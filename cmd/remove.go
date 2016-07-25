// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

const Remove = "remove"

var (
	rmFlagSet *flag.FlagSet
	rmApp     *string
	rmEnv     *string
)

func init() {
	rmFlagSet = flag.NewFlagSet(Remove, flag.ExitOnError)
	rmApp = rmFlagSet.String(AppFlag, "", "app name to be removed")
	rmEnv = rmFlagSet.String(EnvFlag, "", "environment config is for (e.g. PROD, DEV, TEST...)")
}

func RemoveCfg(args []string) error {
	if err := rmFlagSet.Parse(args[2:]); err != nil {
		return err
	}

	if len(*rmEnv) < 1 {
		fmt.Println("NOTE: 'env' flag is not set, defaults to 'default'\n")
		*rmEnv = "default"
	}

	addr := GetConfgrAddr()

	res, err := http.Get(fmt.Sprintf("%s/remove?app=%s&env=%s", addr, *rmApp, *rmEnv))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if code := res.StatusCode; code != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("unable to read remove error response body")
			return err
		}

		fmt.Printf("remove API responded with a status code other than OK: %d\n", code)
		return fmt.Errorf("%s", string(body))
	}

	return nil
}
