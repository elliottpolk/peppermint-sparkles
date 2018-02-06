/// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/main.go
//
package respond

import (
	"fmt"
	"net/http"
)

func WithError(w http.ResponseWriter, statuscode int, format string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}
