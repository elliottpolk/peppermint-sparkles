package respond

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WithJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		fmt.Fprintln(w, `{"success": true}`)
		return
	}

	out, err := json.Marshal(data)
	if err != nil {
		WithError(w, http.StatusInternalServerError, "unable to convert results to json: %v\n", err)
		return
	}

	fmt.Fprintln(w, string(out))
}
