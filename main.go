// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elliottpolk/confgr/config"
	"github.com/elliottpolk/confgr/server"
)

const (
	serverCmd string = "server"
	listCmd   string = "list"
	getCmd    string = "get"
	setCmd    string = "set"
	removeCmd string = "remove"

	appFlag    string = "app"
	configFlag string = "config"
)

func main() {
	getfs := flag.NewFlagSet(getCmd, flag.ExitOnError)
	gaf := getfs.String(appFlag, "", "app name to retrieve respective config")

	setfs := flag.NewFlagSet(setCmd, flag.ExitOnError)
	saf := setfs.String(appFlag, "", "app name to be set")
	scf := setfs.String(configFlag, "", "config to be written")

	removefs := flag.NewFlagSet(removeCmd, flag.ExitOnError)
	raf := removefs.String(appFlag, "", "app name to be removed")

	args := os.Args
	if len(args) == 1 {
		flag.Usage()
		os.Exit(2)
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

			res, err := http.Get(fmt.Sprintf("%s/get?app=%s", serverUrl, *gaf))
			if err != nil {
				log.Fatalf("unable to get config for app %s: %v\n", *gaf, err)
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatalf("unable to retrieve error message: %v\n", err)
			}
			log.Println(string(body))

		case setCmd:
			if err := setfs.Parse(args[2:]); err != nil {
				log.Fatalln(err)
			}

			c := &config.Config{
				App:   *saf,
				Value: *scf,
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

			if res.StatusCode != http.StatusOK {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatalf("unable to retrieve error message: %v\n", err)
				}
				log.Println(string(body))
			}

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

	log.Println("...")

	server.Start()
}
