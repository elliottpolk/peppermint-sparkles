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
	"strings"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/pgp"
	"github.com/elliottpolk/confgr/uuid"
)

const Set = "set"

var (
	setFlagSet  *flag.FlagSet
	setApp      *string
	setEnv      *string
	setCfgValue *string
	setEncrypt  *bool
	setToken    *string
)

func init() {
	setFlagSet = flag.NewFlagSet(Set, flag.ExitOnError)
	setApp = setFlagSet.String(AppFlag, "", "app name to be set")
	setEnv = setFlagSet.String(EnvFlag, "", "environment config is for (e.g. PROD, DEV, TEST...)")
	setCfgValue = setFlagSet.String(CfgFlag, "", "config to be written")
	setEncrypt = setFlagSet.Bool(EncryptFlag, false, "encrypt config")
	setToken = setFlagSet.String(TokenFlag, "", "base64 encoded token to encrypt config with")
}

func SetCfg(args []string) error {
	if err := setFlagSet.Parse(args[2:]); err != nil {
		return err
	}

	if len(*setEnv) < 1 {
		fmt.Println("NOTE: 'env' flag is not set, defaults to 'default'\n")
		*setEnv = "default"
	}

	cfg := &config.Config{
		*setApp,
		*setEnv,
		*setCfgValue,
	}

	var token string
	if *setEncrypt {
		var err error
		token, cfg.Value, err = encryptCfg()
		if err != nil {
			return err
		}
	}

	out, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := sendCfg(string(out)); err != nil {
		return err
	}

	if *setEncrypt {
		fmt.Printf("token: 			 %s\n", token)
		fmt.Printf("token as base64: %s\n", base64.StdEncoding.EncodeToString([]byte(token)))
	}

	fmt.Printf("stored config:\n%s\n", string(out))
	return nil
}

func encryptCfg() (token, cfg string, err error) {
	if len(*setToken) < 1 {
		if token = uuid.GetV4(); len(token) < 1 {
			err = fmt.Errorf("UUID produced an empty string\n")
			return
		}
	} else {
		var t []byte
		if t, err = base64.StdEncoding.DecodeString(*setToken); err != nil {
			return
		}

		token = string(t)
	}

	cypher, err := pgp.Encrypt([]byte(token), []byte(*setCfgValue))
	if err != nil {
		return
	}

	cfg = string(cypher)
	return
}

func sendCfg(cfg string) error {
	addr := GetConfgrAddr()

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/set", addr), strings.NewReader(cfg))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if code := res.StatusCode; code != http.StatusOK {
		fmt.Printf("server responded with status code %d\n", code)
		return fmt.Errorf("%v", body)
	}

	return nil
}
