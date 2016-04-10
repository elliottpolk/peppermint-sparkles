// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/takkun1946/confgr/config"

	"github.com/golang/glog"
)

func init() {
	http.HandleFunc("/list", list)
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)
	http.HandleFunc("/remove", remove)
}

//  list is a http.Handlefunc that retrieves a list of config app names
func list(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, config.ListApps())
}

//  get is a http.HandleFunc that retrieves the relevant config for the provided
//  "app" query parameter. The results are run through json.MarshalIndent prior
//  to writing back to the HTTP reponse.
func get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	app := r.URL.Query().Get("app")
	if len(app) < 1 {
		glog.Errorln("empty app name provided")
		http.Error(w, "invalid app name provided", http.StatusBadRequest)
		return
	}

	out, err := json.MarshalIndent(config.Get(app), "", "   ")
	if err != nil {
		glog.Errorln("unable to marshal conf: %v\n", err)
		http.Error(w, "unable to format data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(out))
}

//  set is a http.HandleFunc that expects a valid config.Config in the request
//  body. It attempts to save the config to the datastore, overwriting the existing
//  app config if one exists. If the save writes properly, a http.StatusOK (200)
//  is returned automatically upon return.
func set(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Errorf("unable to read request body: %v\n", err)
		http.Error(w, "unable to read request", http.StatusInternalServerError)
		return
	}

	c := &config.Config{}
	if err := json.Unmarshal(data, &c); err != nil {
		glog.Errorf("unable to unmarshal content: %v\n", err)
		http.Error(w, "invalid content submitted", http.StatusBadRequest)
		return
	}

	if err := c.Save(); err != nil {
		glog.Errorf("unable to set config %s for app %s: %v\n", err, c.App, c.Value)
		http.Error(w, "unable to add key / value", http.StatusInternalServerError)
		return
	}
}

//  remove is a http.HandleFunc that removes the relevant config for the provided
//  "app" query parameter. If no config exists, no error is returned. If the remove
//  occurs properly, a http.StatusOK (200) is returned automatically upon return.
func remove(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	app := r.URL.Query().Get("app")
	if len(app) < 1 {
		glog.Errorln("empty app name provided")
		http.Error(w, "invalid app name provided", http.StatusBadRequest)
		return
	}

	if err := config.Remove(app); err != nil {
		glog.Errorf("unable to remove config for app %s: %v\n", app, err)
		http.Error(w, "unable to remove app config", http.StatusInternalServerError)
		return
	}
}
