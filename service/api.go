// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package service

import (
	"io/ioutil"
	"net/http"

	"git.platform.manulife.io/oa-montreal/campx/config"
	"git.platform.manulife.io/oa-montreal/campx/log"
	"git.platform.manulife.io/oa-montreal/campx/respond"
)

const (
	PathFind   string = "/api/v1/configs"
	PathSet    string = "/api/v1/set"
	PathRemove string = "/api/v1/remove"

	AppParam string = "app"
	EnvParam string = "env"
)

//  get is a http.HandleFunc that retrieves the relevant config for the provided
//  'app' and 'env' query parameters. An empty 'env' parameter is valid as the
//	config.Get will use 'default'. The results are run through json.MarshalIndent
//	prior to writing back to the HTTP reponse.
func Find(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := r.URL.Query()
	respond.WithJson(w, config.Find(params.Get(AppParam), params.Get(EnvParam)))
}

//  set is a http.HandleFunc that expects a valid config.Config in the request
//  body. It attempts to save the config to the datastore, overwriting the existing
//  app and environment config if one exists. If the save writes properly, a
//	http.StatusOK (200) with the output of what was stored is returned.
func Set(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err, "unable to read request body")
		respond.WithError(w, http.StatusBadRequest, "unable to read request body")

		return
	}

	cfg, err := config.NewConfig(string(data))
	if err != nil {
		log.Error(err, "unable to convert config from string")
		respond.WithError(w, http.StatusBadRequest, "unable to parse config")

		return
	}

	if err := cfg.Save(); err != nil {
		log.Errorf(err, "unable to save config for %s - %s", cfg.App, cfg.Environment)
		respond.WithError(w, http.StatusInternalServerError, "unable to save config")

		return
	}

	respond.WithJson(w, cfg)
}

//  remove is a http.HandleFunc that removes the relevant config for the provided
//  'app' and 'env' query parameters. An empty 'env' parameter is valid as the
//	config.Remove will use 'default'. If no config exists, no error is returned.
//	If the remove occurs properly, a http.StatusOK (200) is returned automatically
//	upon return.
func Remove(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := r.URL.Query()

	app := params.Get(AppParam)
	if len(app) < 1 {
		respond.WithError(w, http.StatusBadRequest, "invalid app name provided")
		return
	}

	env := params.Get(EnvParam)
	if err := config.Remove(app, env); err != nil {
		log.Errorf(err, "unable to remove config for %s - %s", app, env)
		respond.WithError(w, http.StatusInternalServerError, "unable to remove app config")

		return
	}

	respond.WithDefaultOk(w)
}
