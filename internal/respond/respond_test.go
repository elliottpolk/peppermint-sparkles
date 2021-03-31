// Created by Elliott Polk on 28/11/2016
// Copyright © 2016 Manulife AM. All rights reserved.
// internal/respond/respond_test.go
//
package respond

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	go func() {
		if err := http.ListenAndServe(":5000", nil); err != nil {
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

	res, err := http.Get("http://localhost:5000/json")
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
		WithError(w, http.StatusBadRequest, errors.New("fake error"), expected)
	})

	res, err := http.Get("http://localhost:5000/error")
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
