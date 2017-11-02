// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/elliottpolk/confgr/datastore"
	"github.com/elliottpolk/confgr/log"
	"github.com/elliottpolk/confgr/service"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	DefaultStdPort string = "8080"
	DefaultTlsPort string = "8443"

	EnvStdPort string = "CONFGR_STD_PORT"
	EnvTlsPort string = "CONFGR_TLS_PORT"
)

func Serve(context *cli.Context) {
	context.Command.VisibleFlags()

	dsf := context.String(Simplify(DatastoreFlag.Name))
	if len(dsf) < 1 {
		log.NewError("a valid datastore value must be provided")
		return
	}

	//	ensure the expected directory exists
	dir := filepath.Dir(dsf)
	if _, err := os.Stat(dir); err != nil {
		if !os.IsNotExist(err) {
			log.Error(err, "unable to access datastore directory")
			return
		}

		if err := os.MkdirAll(dir, 0766); err != nil {
			log.Error(err, "unable to create datastore directory")
			return
		}
	}

	ds, err := datastore.Open(dsf)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close(false)

	log.Debug("confgr datastore opened")

	http.HandleFunc(service.PathFind, service.Find)
	http.HandleFunc(service.PathSet, service.Set)
	http.HandleFunc(service.PathRemove, service.Remove)

	go func() {
		cert := context.String(TlsCertFlag.Name)
		key := context.String(TlsKeyFlag.Name)

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

		addr := fmt.Sprintf(":%s", DefaultTlsPort)
		if p := os.Getenv(EnvTlsPort); len(p) > 0 {
			addr = fmt.Sprintf(":%s", p)
		}

		log.Debug("starting HTTPS listener")
		log.Fatal(http.ListenAndServeTLS(addr, cert, key, nil))
	}()

	addr := fmt.Sprintf(":%s", DefaultStdPort)
	if p := os.Getenv(EnvStdPort); len(p) > 0 {
		addr = fmt.Sprintf(":%s", p)
	}

	log.Debug("starting HTTP listener")
	log.Fatal(http.ListenAndServe(addr, nil))
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
