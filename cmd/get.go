// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/pgp"
)

const Get = "get"

var (
	getFlagSet *flag.FlagSet
	getApp     *string
	getEnv     *string
	getDecrypt *bool
	getToken   *string
)

func init() {
	getFlagSet = flag.NewFlagSet(Get, flag.ExitOnError)
	getApp = getFlagSet.String(AppFlag, "", "app name to retrieve respective config")
	getEnv = getFlagSet.String(EnvFlag, "", "environment config is for (e.g. PROD, DEV, TEST...)")
	getDecrypt = getFlagSet.Bool(DecryptFlag, false, "decrypt config")
	getToken = getFlagSet.String(TokenFlag, "", "token to decrypt config with")
}

func GetCfg(args []string) error {
	if err := getFlagSet.Parse(args[2:]); err != nil {
		return err
	}

	if *getDecrypt && len(*getToken) < 1 {
		fmt.Println("decryption token must be provided if decryption flag is set to true")
		flag.Usage()
	}

	if len(*getEnv) < 1 {
		fmt.Println("NOTE: 'env' flag is not set, defaults to 'default'\n")
		*getEnv = "default"
	}

	addr := GetConfgrAddr()

	res, err := http.Get(fmt.Sprintf("%s/get?app=%s&env=%s", addr, *getApp, *getEnv))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("\n%s\n", string(body))

	if *getDecrypt {
		return decryptCfg(body)
	}

	return nil
}

func decryptCfg(data []byte) error {
	cfg := &config.Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	t, err := base64.StdEncoding.DecodeString(*getToken)
	if err != nil {
		return err
	}

	plaintxt, err := pgp.Decrypt(t, []byte(cfg.Value))
	if err != nil {
		return err
	}

	fmt.Println("decrypted config:")
	fmt.Printf("%s\n", string(plaintxt))

	return nil
}
