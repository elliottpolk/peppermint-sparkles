/// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package respond

import (
	"fmt"
	"net/http"
)

func WithDefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprint(w, "ok")
}
