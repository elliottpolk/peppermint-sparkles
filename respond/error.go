package respond

import (
	"fmt"
	"net/http"
)

func WithError(w http.ResponseWriter, statuscode int, format string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(format, args...), statuscode)
}
