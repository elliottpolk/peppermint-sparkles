package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const tag string = "peppermint-sparkles.middleware"

func Handle(pattern string, fn http.HandlerFunc) {
	http.Handle(pattern, HandlerFunc(fn))
}

func Handler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request - %+v", r)

		h.ServeHTTP(w, r)
	}
}

func HandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Request - %+v", r)

		fn(w, r)
	}
}
