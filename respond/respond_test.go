//	+build ignore

package respond

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func init() {
	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			fmt.Printf("unable to start http listener: %v\n", err)
			os.Exit(1)
		}
	}()
}

func TestWithJson(t *testing.T) {
	expected := `{"Foo":"Bar"}`
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		WithJson(w, &struct{ Foo string }{"Bar"})
	})

	res, err := http.Get("http://localhost:8000/json")
	if err != nil {
		t.Fatalf("unable to retrieve response from test server: %v\n", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unable to read request body: %v\n", err)
		return
	}

	got := strings.TrimSpace(string(b))
	if got != expected {
		t.Errorf("expected %s - got %s\n", expected, got)
	}
}

func TestWithError(t *testing.T) {
	expected := "error test"
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		WithError(w, http.StatusBadRequest, expected)
	})

	res, err := http.Get("http://localhost:8000/error")
	if err != nil {
		t.Fatalf("unable to retrieve response from test server: %v\n", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unable to read request body: %v\n", err)
		return
	}

	got := strings.TrimSpace(string(b))
	if got != expected {
		t.Errorf("expected %s - got %s\n", expected, got)
	}
}
