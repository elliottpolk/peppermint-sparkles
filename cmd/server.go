// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/cmd/server.go
//
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/campx/backend"
	fileds "git.platform.manulife.io/oa-montreal/campx/backend/file"
	redisds "git.platform.manulife.io/oa-montreal/campx/backend/redis"
	"git.platform.manulife.io/oa-montreal/campx/service"

	bolt "github.com/coreos/bbolt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

//	server flags
var (
	StdListenPortFlag = cli.StringFlag{
		Name:   "p, port",
		Value:  "8080",
		Usage:  "HTTP port to listen on",
		EnvVar: "CAMPX_HTTP_PORT",
	}

	TlsListenPortFlag = cli.StringFlag{
		Name:   "tls-port",
		Value:  "8443",
		Usage:  "HTTPS port to listen on",
		EnvVar: "CAMPX_HTTPS_PORT",
	}

	TlsCertFlag = cli.StringFlag{
		Name:   "tls-cert",
		Usage:  "TLS certificate file for HTTPS",
		EnvVar: "CAMPX_TLS_CERT",
	}

	TlsKeyFlag = cli.StringFlag{
		Name:   "tls-key",
		Usage:  "TLS key file for HTTPS",
		EnvVar: "CAMPX_TLS_KEY",
	}

	DatastoreTypeFlag = cli.StringFlag{
		Name:   "dst, datastore-type",
		Value:  backend.File,
		Usage:  "backend type to be used for storage",
		EnvVar: "CAMPX_DS_TYPE",
	}

	DatastoreFileFlag = cli.StringFlag{
		Name:   "dsf, datastore-file",
		Value:  "/var/lib/confgr/campx.db",
		Usage:  "name / location of file for storing secrets",
		EnvVar: "CAMPX_DS_FILE",
	}

	DatastoreAddrFlag = cli.StringFlag{
		Name:   "dsa, datastore-addr",
		Value:  "localhost:6379",
		Usage:  "address for the remote datastore",
		EnvVar: "CAMPX_DS_ADDR",
	}
)

func Serve(context *cli.Context) {
	context.Command.VisibleFlags()

	var (
		ds  backend.Datastore
		err error
	)

	dst := context.String(flag(DatastoreTypeFlag.Name))

	switch dst {
	case backend.Redis:

		//	TODO ... include / handle additional Redis options (e.g. timeout, etc)
		addr := context.String(flag(DatastoreAddrFlag.Name))
		if ds, err = redisds.Open(&redis.Options{Addr: addr}); err != nil {
			log.Error(err, "unable to open connection to datastore")
			return
		}

	case backend.File:

		//	TODO ... include / handle additional bolt options (e.g. timeout, etc)
		fname := context.String(flag(DatastoreFileFlag.Name))
		if ds, err = fileds.Open(fname, bolt.DefaultOptions); err != nil {
			log.Error(err, "unable to open connection to datastore")
			return
		}

	default:
		log.Errorf("%s is not a supported datastore type", dst)
		return
	}

	defer ds.Close()
	log.Debug("datastore opened")

	mux := http.NewServeMux()

	//	attach current service handler
	mux = service.Handle(mux, &service.Handler{Backend: ds})

	//	start HTTPS listener in a seperate go routine since it is a blocking func
	go func() {
		cert, key := context.String(flag(TlsCertFlag.Name)), context.String(flag(TlsKeyFlag.Name))
		if len(cert) < 1 || len(key) < 1 {
			return
		}

		if _, err := os.Stat(cert); err != nil {
			log.Error(err, "unable to access TLS cert file")
			return
		}

		if _, err := os.Stat(key); err != nil {
			log.Error(err, "unable to access TLS key file")
			return
		}

		addr := fmt.Sprintf(":%s", context.String(flag(TlsListenPortFlag.Name)))

		log.Debug("starting HTTPS listener")
		log.Fatal(http.ListenAndServeTLS(addr, cert, key, mux))
	}()

	log.Debug("starting HTTP listener")

	addr := fmt.Sprintf(":%s", context.String(flag(StdListenPortFlag.Name)))
	log.Fatal(http.ListenAndServe(addr, mux))
}

func asURL(addr, path, params string) string {
	scheme := "http"
	if https := "https"; strings.HasPrefix(addr, https) {
		addr = strings.TrimPrefix(addr, fmt.Sprintf("%s://", https))
		scheme = https
	}

	//	attempt to scrub the scheme if it was not handled above
	addr = strings.TrimPrefix(addr, "http://")

	return (&url.URL{
		Scheme:   scheme,
		Host:     addr,
		Path:     path,
		RawQuery: params,
	}).String()
}

func retrieve(from string) (string, error) {
	res, err := http.Get(from)
	if err != nil {
		return "", errors.Wrap(err, "unable to call service")
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable to read service response body")
	}

	if code := res.StatusCode; code != http.StatusOK {
		return "", errors.Errorf("confgr service responded with status code %d and message %s", code, string(b))
	}

	return string(b), nil
}

func send(to, body string) (string, error) {
	res, err := http.Post(to, http.DetectContentType([]byte(body)), strings.NewReader(body))
	if err != nil {
		return "", errors.Wrap(err, "unable to post config")
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "unable read service response body")
	}

	if code := res.StatusCode; code != http.StatusOK {
		return "", errors.Errorf("confgr service responded with status code %d and message %s", code, string(b))
	}

	return string(b), nil
}
