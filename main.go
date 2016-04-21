// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/pgp"
	"github.com/elliottpolk/confgr/server"
	"github.com/elliottpolk/confgr/uuid"
)

const (
	serverCmd string = "server"
	listCmd   string = "list"
	getCmd    string = "get"
	setCmd    string = "set"
	removeCmd string = "remove"

	appFlag     string = "app"
	configFlag  string = "config"
	encryptFlag string = "encrypt"
	decryptFlag string = "decrypt"
	tokenFlag   string = "token"
)

func main() {
	flag.Lookup("alsologtostderr").Value.Set("true")
	flag.Usage = usage

	getfs := flag.NewFlagSet(getCmd, flag.ExitOnError)
	gaf := getfs.String(appFlag, "", "app name to retrieve respective config")
	gdf := getfs.Bool(decryptFlag, false, "decrypt config")
	gtf := getfs.String(tokenFlag, "", "token to decrypt config with")

	setfs := flag.NewFlagSet(setCmd, flag.ExitOnError)
	saf := setfs.String(appFlag, "", "app name to be set")
	scf := setfs.String(configFlag, "", "config to be written")
	sef := setfs.Bool(encryptFlag, false, "encrypt config")

	removefs := flag.NewFlagSet(removeCmd, flag.ExitOnError)
	raf := removefs.String(appFlag, "", "app name to be removed")

	args := os.Args
	if len(args) == 1 {
		flag.Usage()
	}

	for _, a := range args[1:] {
		if a == "-h" || a == "-help" || a == "--help" {
			flag.Usage()
		}
	}

	if args[1] != serverCmd {
		if os.Getenv("CONFGR_ADDR") == "" {
			log.Fatalln("CONFGR_ADDR must be set prior to usage (i.e. export CONFGR_ADDR=localhost:8080)")
		}

		serverUrl := os.Getenv("CONFGR_ADDR")
		if !strings.HasPrefix(serverUrl, "http") {
			serverUrl = fmt.Sprintf("http://%s", serverUrl)
		}

		if strings.HasSuffix(serverUrl, "/") {
			serverUrl = strings.TrimSuffix(serverUrl, "/")
		}

		switch args[1] {
		case listCmd:
			res, err := http.Get(fmt.Sprintf("%s/list", serverUrl))
			if err != nil {
				log.Printf("unable to list apps: %v\n", err)
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Printf("unable to list apps: %v\n", err)
				return
			}
			log.Println(string(body))

		case getCmd:
			if err := getfs.Parse(args[2:]); err != nil {
				log.Fatalln(err)
			}

			if *gdf && len(*gtf) < 1 {
				log.Println("decryption must be provided if decryption is specified")
				flag.Usage()
			}

			res, err := http.Get(fmt.Sprintf("%s/get?app=%s", serverUrl, *gaf))
			if err != nil {
				log.Fatalf("unable to get config for app %s: %v\n", *gaf, err)
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("unable to retrieve error message: %v\n", err)
			}
			log.Printf("\n%s\n", string(body))

			if *gdf {
				t, err := base64.StdEncoding.DecodeString(*gtf)
				if err != nil {
					log.Fatalf("unable to convert base64 token to string: %v\n", err)
				}

				cfg := &config.Config{}
				if err := json.Unmarshal(body, &cfg); err != nil {
					log.Fatalf("unable to convert response body to Config: %v\n", err)
				}

				ctxt, err := base64.StdEncoding.DecodeString(cfg.Value)
				if err != nil {
					log.Fatalf("unable to base64 decode config value: %v\n", err)
				}

				ptxt, err := pgp.Decrypt(t, ctxt)
				if err != nil {
					log.Fatalf("unable to decrypt config: %v\n", err)
				}

				log.Printf("\n%s\n", string(ptxt))
			}

		case setCmd:
			if err := setfs.Parse(args[2:]); err != nil {
				log.Fatalln(err)
			}

			c := &config.Config{
				App:   *saf,
				Value: *scf,
			}

			var token string
			if *sef {
				if token = uuid.GetV4(); len(token) < 1 {
					log.Fatalf("encryption token produced an empty string\n")
				}

				val, err := pgp.Encrypt([]byte(token), []byte(*scf))
				if err != nil {
					log.Fatalf("unable to encrypt config: %v\n", err)
				}

				c.Value = string(val)
			}

			out, err := json.Marshal(c)
			if err != nil {
				log.Fatalf("unable to prepare request: %v\n", err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/set", serverUrl), strings.NewReader(string(out)))
			if err != nil {
				log.Fatalf("unable to prepare request: %v\n", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalf("unable to set config for app %s: %v\n", c.App, err)
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("unable to retrieve response: %v\n", err)
			}

			if *sef {
				log.Printf("token: 			 %s", token)
				log.Printf("token as base64: %s", base64.StdEncoding.EncodeToString([]byte(token)))
			}
			log.Printf("stored config:\n%s\n", string(body))

		case removeCmd:
			if err := removefs.Parse(args[2:]); err != nil {
				log.Fatalln(err)
			}

			res, err := http.Get(fmt.Sprintf("%s/remove?app=%s", serverUrl, *raf))
			if err != nil {
				log.Fatalf("unable to remove config for app %s: %v\n", *raf, err)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatalf("unable to retrieve error message: %v\n", err)
				}
				log.Println(string(body))
			}

		default:
			flag.Usage()
		}

		return
	}

	server.Start()
}

func usage() {
	fmt.Printf("usage: %s <command> [args]\n\n", os.Args[0])

	fmt.Println("Available commands:")
	fmt.Printf("\t%s\t\tstarts confgr server\n", serverCmd)
	fmt.Printf("\t%s\t\tretrieves the available app configs\n", listCmd)
	fmt.Printf("\t%s\t\tretrieves the available config\n", getCmd)
	fmt.Printf("\t%s\t\tadds a new config\n", setCmd)
	fmt.Printf("\t%s\t\tdeletes the specified config\n", removeCmd)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case getCmd:
			fmt.Printf("Arguments for %s:\n", getCmd)
			fmt.Printf("\t%s\tapp name to retrieve respective config\n", appFlag)
			fmt.Printf("\t%s\tdecrypt the config\n", decryptFlag)
			fmt.Printf("\t%s\ttoken to decrypt config with\n", tokenFlag)

		case setCmd:
			fmt.Printf("Arguments for %s:\n", setCmd)
			fmt.Printf("\t%s\tapp name to be set\n", appFlag)
			fmt.Printf("\t%s\tconfig to be written\n", configFlag)
			fmt.Printf("\t%s\tencrypt config\n", encryptFlag)

		case removeCmd:
			fmt.Printf("Arguments for %s:\n", removeCmd)
			fmt.Printf("\t%s\tapp name to be removed\n", appFlag)
		}
	}

	fmt.Println()
	os.Exit(0)
}
