// Created by Elliott Polk on 28/11/2016
// Copyright © 2016 Manulife AM. All rights reserved.
// internal/respond/json.go
//
package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintln(w, `{}`)
		return
	}

	out, err := json.Marshal(data)
	if err != nil {
		WithError(w, http.StatusInternalServerError, err, "unable to convert results to json")
		return
	}

	fmt.Fprintln(w, string(out))
}

func WithJsonCreated(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if data != nil {
		out, err := json.Marshal(data)
		if err != nil {
			WithError(w, http.StatusInternalServerError, err, "unable to convert results to json")
			return
		}

		fmt.Fprintln(w, string(out))
	}
}
