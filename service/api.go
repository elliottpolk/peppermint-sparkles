package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/go-common/respond"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/middleware"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"

	"github.com/pkg/errors"
)

const tag string = "peppermint-sparkles.service"

const (
	PathSecrets string = "/api/v2/secrets"

	AppParam  string = "app_name"
	EnvParam  string = "env"
	UserParam string = "username"
	IdParam   string = "uuid"
)

var idExp *regexp.Regexp = regexp.MustCompile(`secrets/(?P<id>([a-zA-Z\d]+(-)?){5})(\/)?$`)

type Handler struct {
	Backend backend.Datastore
}

func Handle(mux *http.ServeMux, h *Handler) *http.ServeMux {
	mux.Handle(PathSecrets, middleware.Handler(h))
	mux.Handle(fmt.Sprintf("%s/", PathSecrets), middleware.Handler(h))
	return mux
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	matched, id, err := getId(r.URL.Path)
	if err != nil {
		log.Error(tag, err, "unable to retrieve the secret ID from the URL path")
		respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
		return
	}

	ds, params := h.Backend, r.URL.Query()

	if !matched {
		//	if neither the ID nor app name and environment combination are provided,
		//	there is really no way to retrieve a secret
		app, env := params.Get(AppParam), params.Get(EnvParam)
		if len(app) < 1 || len(env) < 1 {
			respond.WithErrorMessage(w, http.StatusBadRequest, "a valid ID or %s and %s must be specified", AppParam, EnvParam)
			return
		}

		for _, k := range ds.Keys() {
			if strings.HasSuffix(k, backend.KeySuffix(app, env)) {
				id = k
				break
			}
		}
		log.Debugf(tag, "attempted to find an ID for app %s and env %s: %s", app, env, id)
	}

	raw := ds.Get(id)
	if len(raw) < 1 {
		respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
		return
	}

	rec, err := models.ParseRecord(raw)
	if err != nil {
		log.Error(tag, err, "unable to parse stored secret")
		respond.WithErrorMessage(w, http.StatusBadRequest, "invalid secret")
		return
	}

	if rec.Status != models.ActiveStatus {
		log.Infof(tag, "record for ID %s found, but has status %s", rec.Id, rec.Status)
		respond.WithErrorMessage(w, http.StatusNotFound, "file not found")
		return
	}

	log.Debugf(tag, "retrieved secret with ID %s", id)
	respond.WithJson(w, rec.Secret)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ds, params := h.Backend, r.URL.Query()

	in, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(tag, err, "unable to read in request body")
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
		log.Error(tag, err, "unable to unmarshal request to secret")
		respond.WithErrorMessage(w, http.StatusBadRequest, "unable to convert request to valid secret")
		return
	}

	if len(s.Id) < 1 {
		respond.WithErrorMessage(w, http.StatusBadRequest, "an ID for the secret must be specified")
		return
	}

	now := time.Now().UnixNano()
	rec := &models.Record{
		Secret:    s,
		Created:   now,
		CreatedBy: usr,
		Updated:   now,
		UpdatedBy: usr,
		Status:    models.ActiveStatus,
	}

	if rec.Exists(ds) {
		respond.WithErrorMessage(w, http.StatusConflict, "record found for provided ID")
		return
	}

	if err := rec.Write(ds); err != nil {
		log.Error(tag, err, "unable to write record to storage")
		respond.WithErrorMessage(w, http.StatusInternalServerError, "unable to write secret record to storage")
		return
	}

	log.Debugf(tag, "created new record with ID %s for user %s", s.Id, usr)
	respond.WithJsonCreated(w, s)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.get(w, r)

	case http.MethodPost:
		h.create(w, r)

	default:
		respond.WithMethodNotAllowed(w)
		return
	}
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
