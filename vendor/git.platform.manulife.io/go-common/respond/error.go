// Created by Elliott Polk on 28/11/2016
// Copyright © 2016 Manulife AM. All rights reserved.
// go-common/respond/error.go
//
package respond

import (
	"fmt"
	"net/http"

	"git.platform.manulife.io/go-common/log"
)

func WithErrorMessage(w http.ResponseWriter, statuscode int, format string, args ...interface{}) {
	log.Errorf(format, args...)
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}

func WithNewError(w http.ResponseWriter, statuscode int, message string) {
	log.Errorf("%s", message)
	http.Error(w, message, statuscode)
}

func WithErrorf(w http.ResponseWriter, statuscode int, err error, format string, args ...interface{}) {
	log.Error(err, fmt.Sprintf(format, args...))
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}

func WithError(w http.ResponseWriter, statuscode int, err error, message string) {
	log.Error(err, message)
	http.Error(w, message, statuscode)
}

func WithMethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}
