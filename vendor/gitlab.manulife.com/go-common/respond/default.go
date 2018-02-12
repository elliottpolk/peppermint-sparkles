// Created by Elliott Polk on 17/03/2017
// Copyright © 2017 Manulife AM. All rights reserved.
// go-common/respond/default.go
//
package respond

import (
	"fmt"
	"net/http"
)

func WithDefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprintln(w, "ok")
}
