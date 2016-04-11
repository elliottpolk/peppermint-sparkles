// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/takkun1946/confgr/config"
	"github.com/takkun1946/confgr/server"

	"github.com/golang/glog"
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
	flag.Usage = func() {
		fmt.Println("usage: confgr command [arguments]\n")
		fmt.Println("commands:")
		fmt.Printf("\t%s\tstarts the confgr server\n", serverCmd)
		fmt.Printf("\t%s\tlists out the available app configs\n", listCmd)
		fmt.Printf("\t%s\tretrieves the config for the provided app name\n", getCmd)
		fmt.Printf("\t%s\tsets the config for the specified app\n", setCmd)
		fmt.Printf("\t%s\tremoves the config and app for the provided app\n", removeCmd)

		switch flag.Arg(0) {
		case getCmd:
			fmt.Printf("arguments for %s:\n", getCmd)
			fmt.Printf("\t%s\tapp name to retrieve respective config\n", appFlag)

		case setCmd:
			fmt.Printf("arguments for %s:\n", setCmd)
			fmt.Printf("\t%s\tapp name to be set\n", appFlag)
			fmt.Printf("\t%s\tconfig to be written\n", configFlag)

		case removeCmd:
			fmt.Printf("arguments for %s:\n", removeCmd)
			fmt.Printf("\t%s\tapp name to be removed\n", appFlag)

		}

		os.Exit(0)
	}

	flag.Parse()
	defer glog.Flush()

	//  force log output to stdout / stderr
	flag.Lookup("alsologtostderr").Value.Set("true")

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
	}

	for _, a := range args[1:] {
		if a == "-h" || a == "-help" || a == "--help" {
			flag.Usage()
		}
	}

	if args[0] != serverCmd {
		if os.Getenv("CONFGR_ADDR") == "" {
			fmt.Println("CONFGR_ADDR must be set prior to usage (i.e. export CONFGR_ADDR=localhost:8080)")
			os.Exit(1)
		}

		serverUrl := os.Getenv("CONFGR_ADDR")
		if !strings.HasPrefix(serverUrl, "http") {
			serverUrl = fmt.Sprintf("http://%s", serverUrl)
		}

		if strings.HasSuffix(serverUrl, "/") {
			serverUrl = strings.TrimSuffix(serverUrl, "/")
		}

		switch args[0] {
		case listCmd:
			res, err := http.Get(fmt.Sprintf("%s/list", serverUrl))
			if err != nil {
				glog.Errorf("unable to list apps: %v\n", err)
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				glog.Errorf("unable to list apps: %v\n", err)
				return
			}
			fmt.Println(string(body))

		case getCmd:
			if len(flag.Args()) < 2 {
				flag.Usage()
			}

			appName := flag.Arg(1)

			res, err := http.Get(fmt.Sprintf("%s/get?app=%s", serverUrl, appName))
			if err != nil {
				glog.Errorf("unable to get config for app %s: %v\n", appName, err)
				return
			}
			defer res.Body.Close()

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				glog.Errorf("unable to retrieve error message: %v\n", err)
				return
			}
			fmt.Println(string(body))

		case setCmd:
			if len(flag.Args()) < 3 {
				flag.Usage()
			}

			c := &config.Config{
				App:   flag.Arg(1),
				Value: flag.Arg(2),
			}

			out, err := json.Marshal(c)
			if err != nil {
				glog.Errorf("unable to prepare request: %v\n", err)
				return
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/set", serverUrl), strings.NewReader(string(out)))
			if err != nil {
				glog.Errorf("unable to prepare request: %v\n", err)
				return
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				glog.Errorf("unable to set config for app %s: %v\n", c.App, err)
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					glog.Errorf("unable to retrieve error message: %v\n", err)
					return
				}
				fmt.Println(string(body))
			}

		case removeCmd:
			if len(flag.Args()) < 2 {
				flag.Usage()
			}

			appName := flag.Arg(1)
			res, err := http.Get(fmt.Sprintf("%s/remove?app=%s", serverUrl, appName))
			if err != nil {
				glog.Errorf("unable to remove config for app %s: %v\n", appName, err)
				return
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					glog.Errorf("unable to retrieve error message: %v\n", err)
					return
				}
				fmt.Println(string(body))
			}

		default:
			flag.Usage()
		}

		return
	}

	server.Start()
}
