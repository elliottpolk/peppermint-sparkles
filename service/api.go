// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/service/api.go
//
package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"gitlab.manulife.com/go-common/log"
	"gitlab.manulife.com/go-common/respond"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/backend"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/middleware"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/secret"

	"github.com/pkg/errors"
)

const (
	PathSecrets string = "/api/v1/secrets"

	AppParam string = "app_name"
	EnvParam string = "env"
)

var idExp *regexp.Regexp = regexp.MustCompile(`secrets/(?P<id>[a-zA-Z\d]+)(\/)?$`)

type Handler struct {
	Backend backend.Datastore
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	matched, id, err := getId(r.URL.Path)
	if err != nil {
		log.Error(err, "unable to parse ")
		respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
		return
	}

	//	pre-checks for paths without a valid ID
	if !matched {
		switch r.Method {
		case http.MethodPut, http.MethodDelete:
			respond.WithMethodNotAllowed(w)
			return
		case http.MethodGet:
			//	if the ID nor app name is not provided, there is really no way to
			//	retrieve a secret
			if app := r.URL.Query().Get(AppParam); len(app) < 1 {
				respond.WithErrorMessage(w, http.StatusBadRequest, "a valid app must be specified")
				return
			}
		}
	}

	ds := h.Backend

	switch r.Method {
	case http.MethodGet:
		for !matched && len(id) < 1 {
			params := r.URL.Query()

			app := params.Get(AppParam)
			//	if an environment was provided, convert app + env value to a backend
			//	key and attempt to retrieve
			if env := params.Get(EnvParam); len(env) > 0 {
				id = backend.Key(app, env)
				break
			}

			//	no environment value was provided, so list out all apps and attempt
			//	to filter out the secrets for the provided app name
			res, err := ds.List()
			if err != nil {
				log.Error(err, "unable to list out current values from datastore")
				respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to retrieve secrets")
				return
			}

			secrets := make([]*secret.Secret, 0)
			for _, r := range res {
				for _, v := range r {
					s, err := secret.NewSecret(v)
					if err != nil {
						//	this is likely a larger issue, so make no assumptions and bail
						log.Error(err, "unable to parse raw string to secret")
						respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to retrieve secrets")
						return
					}

					if s.App == app {
						secrets = append(secrets, s)
					}
				}
			}

			//	return turn filtered results
			respond.WithJson(w, secrets)
			return
		}

		raw := ds.Get(id)
		if len(raw) < 1 {
			respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
			return
		}

		s, err := secret.NewSecret(raw)
		if err != nil {
			log.Error(err, "unable to parse stored secret")
			respond.WithError(w, http.StatusBadRequest, err, "invalid secret")
			return
		}

		//	this means that the URL had the ID in place and only 1 secret should
		//	returned
		if matched {
			respond.WithJson(w, s)
			return
		}

		//	this means that the URL did not have the ID in place and there is
		//	the likelyhood that more than 1 can be returned in other code
		respond.WithJson(w, []*secret.Secret{s})
		return

	case http.MethodPost,
		http.MethodPut:

		in, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err, "unable to read in request body")
			respond.WithError(w, http.StatusBadRequest, err, "unable to read in request")
			return
		}

		//	convert request body to secret to (sort of) ensure data received is
		//	in the expected secret format
		s := &secret.Secret{}
		if err := json.Unmarshal(in, &s); err != nil {
			log.Error(err, "unable to unmarshal request to secret")
			respond.WithError(w, http.StatusBadRequest, err, "unable to convert request to valid secret")
			return
		}

		//	if no id is provided, generate a new "key" using the app name and env
		if !matched {
			id = backend.Key(s.App, s.Env)
		}

		if len(s.Id) < 1 {
			s.Id = id
		}

		//	check if the secret with the current id exists in the datastore
		exists := (len(ds.Get(id)) > 0)

		//	convert back to string before storage
		out, err := json.Marshal(s)
		if err != nil {
			log.Error(err, "unable to marshal secret")
			respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to prep secret for storage")
			return
		}

		if err := ds.Set(id, string(out)); err != nil {
			log.Error(err, "unable to write secret to storage")
			respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to write secret to storage")
			return
		}

		//	respond with 201 if the resource did not exist before
		if !exists {
			respond.WithJsonCreated(w, s)
			return
		}

		respond.WithJson(w, s)
		return

	case http.MethodDelete:
		if err := ds.Remove(id); err != nil {
			log.Errorf("%v: unable to remove secret for id %s", err, id)
			respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to remove secret")
			return
		}

		respond.WithJson(w, nil)
		return

	default:
		respond.WithMethodNotAllowed(w)
		return
	}
}

func Handle(mux *http.ServeMux, h *Handler) *http.ServeMux {
	mux.Handle(PathSecrets, middleware.Handler(h))
	mux.Handle(fmt.Sprintf("%s/", PathSecrets), middleware.Handler(h))
	return mux
}

func getId(path string) (bool, string, error) {
	matched, err := regexp.Match(idExp.String(), []byte(path))
	if err != nil {
		return false, "", errors.Wrap(err, "unable to process path")
	}

	if matched {
		matches, m := idExp.FindStringSubmatch(path), make(map[string]string)
		for i, n := range idExp.SubexpNames() {
			if i > 0 && i <= len(matches) {
				m[n] = matches[i]
			}
		}

		return true, m["id"], nil
	}
	return false, "", nil
}
