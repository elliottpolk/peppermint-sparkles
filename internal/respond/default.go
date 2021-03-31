// Created by Elliott Polk on 17/03/2017
// Copyright © 2017 Manulife AM. All rights reserved.
// internal/respond/default.go
//
package respond

import (
	"fmt"
	"net/http"
)

const tag string = "go-common.respond"

func WithDefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprintln(w, "ok")
}
