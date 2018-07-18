// Created by Elliott Polk on 25/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/middleware/handle.go
//
package middleware

import (
	"net/http"

	"git.platform.manulife.io/go-common/log"
)

const tag string = "manulife.oa-montreal.peppermint-sparkles.middleware"

func Handle(pattern string, fn http.HandlerFunc) {
	http.Handle(pattern, HandlerFunc(fn))
}

func Handler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof(tag, "Request - %+v", r)

		h.ServeHTTP(w, r)
	}
}

func HandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof(tag, "Request - %+v", r)

		fn(w, r)
	}
}
