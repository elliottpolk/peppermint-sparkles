package respond

import (
	"fmt"
	"net/http"
)

func WithDefaultOk(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=us-ascii")
	fmt.Fprint(w, "ok")
}
