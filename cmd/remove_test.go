package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
)

func TestRm(t *testing.T) {
	repo := fmt.Sprintf("test_%s.db", uuid.New().String())
	ds, err := fileds.Open(repo, bolt.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		ds.Close()
		if err := os.RemoveAll(repo); err != nil {
			t.Errorf("unable to remove temporary test repo %s\n", repo)
		}
	}()

	id, app, env, content, usr := uuid.New().String(), "dummy", "test", "notSuperS3cret", "tester"

	sample := fmt.Sprintf(`{"id":"%s","app_name":"%s","env":"%s","content":"%s"}`, id, app, env, content)
	src, err := models.ParseSecret(sample)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UnixNano()
	rec := &models.Record{
		Secret:    src,
		Created:   now,
		CreatedBy: usr,
		Updated:   now,
		UpdatedBy: usr,
		Status:    models.ActiveStatus,
	}

	if err := rec.Write(ds); err != nil {
		t.Fatal(err)
	}

	port := freeport()

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = service.Handle(mux, &service.Handler{Backend: ds})

		wg.Done()
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
	}(ds)
	wg.Wait()

	params := &url.Values{
		service.AppParam:  []string{app},
		service.EnvParam:  []string{env},
		service.UserParam: []string{usr},
	}

	if err := rm(false, id, fmt.Sprintf("http://localhost:%d", port), params); err != nil {
		t.Fatal(err)
	}
}

func TestRmInsecure(t *testing.T) {
	repo := fmt.Sprintf("test_%s.db", uuid.New().String())
	ds, err := fileds.Open(repo, bolt.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		ds.Close()
		if err := os.RemoveAll(repo); err != nil {
			t.Errorf("unable to remove temporary test repo %s\n", repo)
		}
	}()

	id, app, env, content, usr := uuid.New().String(), "dummy", "test", "notSuperS3cret", "tester"

	sample := fmt.Sprintf(`{"id":"%s","app_name":"%s","env":"%s","content":"%s"}`, id, app, env, content)
	src, err := models.ParseSecret(sample)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UnixNano()
	rec := &models.Record{
		Secret:    src,
		Created:   now,
		CreatedBy: usr,
		Updated:   now,
		UpdatedBy: usr,
		Status:    models.ActiveStatus,
	}

	if err := rec.Write(ds); err != nil {
		t.Fatal(err)
	}

	port := freeport()

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = service.Handle(mux, &service.Handler{Backend: ds})

		cert := "testdata/cert.pem"
		key := "testdata/key.pem"

		wg.Done()
		t.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", port), cert, key, mux))
	}(ds)
	wg.Wait()

	params := &url.Values{
		service.AppParam:  []string{app},
		service.EnvParam:  []string{env},
		service.UserParam: []string{usr},
	}

	if err := rm(false, id, fmt.Sprintf("https://localhost:%d", port), params); err != nil && !strings.HasSuffix(err.Error(), "x509: certificate signed by unknown authority") {
		t.Fatal(err)
	}

	if err := rm(true, id, fmt.Sprintf("https://localhost:%d", port), params); err != nil {
		t.Fatal(err)
	}
}
