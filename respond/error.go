/// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package respond

import (
	"fmt"
	"net/http"
)

func WithError(w http.ResponseWriter, statuscode int, format string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}
