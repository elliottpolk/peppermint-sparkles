/// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package respond

import (
	"fmt"
	"net/http"
)

func WithDefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprint(w, "ok")
}
