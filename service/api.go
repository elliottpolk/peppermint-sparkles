// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/service/api.go
//
package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"gitlab.manulife.com/go-common/log"
	"gitlab.manulife.com/go-common/respond"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/backend"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/middleware"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/models"

	"github.com/pkg/errors"
)

const (
	PathSecrets string = "/api/v2/secrets"

	AppParam  string = "app_name"
	EnvParam  string = "env"
	UserParam string = "username"
)

var idExp *regexp.Regexp = regexp.MustCompile(`secrets/(?P<id>[a-zA-Z\d]+)(\/)?$`)

type Handler struct {
	Backend backend.Datastore
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	matched, id, err := getId(r.URL.Path)
	if err != nil {
		log.Error(err, "unable to retrieve the secret ID from the URL path")
		respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
		return
	}

	ds, params := h.Backend, r.URL.Query()

	//	pre-checks for paths without a valid ID
	if !matched {
		switch r.Method {
		case http.MethodPut, http.MethodDelete:
			respond.WithMethodNotAllowed(w)
			return
		case http.MethodGet:
			//	if neither the ID nor app name and environment combination are
			//	provided, there is really no way to retrieve a secret
			if app, env := params.Get(AppParam), params.Get(EnvParam); len(app) < 1 || len(env) < 1 {
				respond.WithErrorMessage(w, http.StatusBadRequest, "a valid %s and %s must be specified", AppParam, EnvParam)
				return
			}
		}
	}

	switch r.Method {
	case http.MethodGet:
		if !matched {
			app, env := params.Get(AppParam), params.Get(EnvParam)
			for _, k := range ds.Keys() {
				if strings.HasSuffix(k, backend.KeySuffix(app, env)) {
					id = k
					break
				}
			}
			log.Debugf("attempted to find an ID for app %s and env %s: %s", app, env, id)
		}

		raw := ds.Get(id)
		if len(raw) < 1 {
			respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
			return
		}

		rec, err := models.ParseRecord(raw)
		if err != nil {
			log.Error(err, "unable to parse stored secret")
			respond.WithErrorMessage(w, http.StatusBadRequest, "invalid secret")
			return
		}

		if rec.Status != models.ActiveStatus {
			log.Infof("record for ID %s found, but has status %s", rec.Id, rec.Status)
			respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
			return
		}
		log.Debugf("retrieved secret with ID %s", id)

		respond.WithJson(w, rec.Secret)
		return

	case http.MethodPost,
		http.MethodPut:

		in, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err, "unable to read in request body")
			respond.WithErrorMessage(w, http.StatusBadRequest, "unable to read in request")
			return
		}

		usr := params.Get(UserParam)
		if len(usr) < 1 {
			respond.WithErrorMessage(w, http.StatusBadRequest, "a valid user name must be provided")
			return
		}

		//	convert request body to secret to (sort of) ensure data received is
		//	in the expected secret format
		s, err := models.ParseSecret(string(in))
		if err != nil {
			log.Error(err, "unable to unmarshal request to secret")
			respond.WithErrorMessage(w, http.StatusBadRequest, "unable to convert request to valid secret")
			return
		}

		now := time.Now().UnixNano()

		if !matched {
			app, env := s.App, s.Env
			for _, k := range ds.Keys() {
				if strings.HasSuffix(k, backend.KeySuffix(app, env)) {
					id = k
					break
				}
			}
			log.Debugf("attempted to find an ID for app %s and env %s: %s", app, env, id)

			//	no record can be found with the data provided, so attempt to
			//	generate a new one
			if len(id) < 1 {
				//	trigger a creation record to start the audit trail for a record
				history := &models.Historical{}
				if err := history.Write(ds, models.CreateAction, usr, now); err != nil {
					log.Error(err, "unable to write historical record")
					respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to write record")
					return
				}
				log.Debugf("wrote historical creation record for user %s", usr)

				//	add in pseudorandom noise along with the app name and env to
				//	attempt to prevent ID collisions
				s.Id = s.NewId()
				rec := &models.Record{
					Secret:    s,
					Created:   now,
					CreatedBy: usr,
					Updated:   now,
					UpdatedBy: usr,
					Status:    models.ActiveStatus,
				}

				if err := rec.Write(ds); err != nil {
					log.Error(err, "unable to store record")
					respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to write secret record to storage")
					return
				}
				log.Debugf("created new record with ID %s for user %s", s.Id, usr)

				respond.WithJsonCreated(w, s)
				return
			}
		}

		//	the ID may not be provided in the submitted secret, so ensure ID is
		//	set from the request.URL.Path or found via the app_name / env combo
		if len(s.Id) < 1 {
			s.Id = id
		}

		//	attempt to retrieve the current record based on the provided secret.Id
		if curr := ds.Get(s.Id); len(curr) > 0 {
			history, err := models.FromCurrent(curr)
			if err != nil {
				log.Error(err, "unable to generate historical record")
				respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to update existing record")
				return
			}

			if err := history.Write(ds, models.UpdateAction, usr, now); err != nil {
				log.Error(err, "unable to write historical record")
				respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to update existing record")
				return
			}
			log.Debugf("wrote historical record with ID %s for user %s with action %s", s.Id, usr, models.UpdateAction)

			rec := &models.Record{
				Secret:    s,
				Created:   history.Created,
				CreatedBy: history.CreatedBy,
				Updated:   now,
				UpdatedBy: usr,
				Status:    models.ActiveStatus,
			}

			if err := rec.Write(ds); err != nil {
				log.Error(err, "unable to write record")
				respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to update existing record")
				return
			}
			log.Debugf("wrote updated record with ID %s for user %s", s.Id, usr)

			respond.WithJson(w, s)
			return
		}

		respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
		return

	case http.MethodDelete:
		usr := params.Get(UserParam)
		if len(usr) < 1 {
			respond.WithErrorMessage(w, http.StatusBadRequest, "a valid user name must be provided")
			return
		}

		now := time.Now().UnixNano()
		if curr := ds.Get(id); len(curr) > 0 {
			history, err := models.FromCurrent(curr)
			if err != nil {
				log.Error(err, "unable to generate historical record")
				respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to update existing record")
				return
			}

			if err := history.Write(ds, models.ArchiveAction, usr, now); err != nil {
				log.Error(err, "unable to write historical record")
				respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to update existing record")
				return
			}
			log.Debugf("wrote historical record with ID %s for user %s with action %s", id, usr, models.ArchiveAction)
		}

		if err := ds.Remove(id); err != nil {
			log.Errorf("%v: unable to remove secret for id %s", err, id)
			respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to remove secret")
			return
		}
		log.Debugf("removed record with ID %s for user %s", id, usr)

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
