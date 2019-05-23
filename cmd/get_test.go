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
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
)

func TestGet(t *testing.T) {

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

	tok, err := crypto.NewToken()
	if err != nil {
		t.Fatal(err)

	}

	id, app, env, content := uuid.New().String(), "dummy", "test", "notSuperS3cret"

	crypter := &pgp.Crypter{Token: []byte(tok)}
	cypher, err := crypter.Encrypt([]byte(content))
	if err != nil {
		t.Fatal(err)
	}

	sample := fmt.Sprintf(`{"id":"%s","app_name":"%s","env":"%s","content":"%s"}`, id, app, env, string(cypher))
	src, err := models.ParseSecret(sample)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UnixNano()
	rec := &models.Record{
		Secret:    src,
		Created:   now,
		CreatedBy: "tester",
		Updated:   now,
		UpdatedBy: "tester",
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
		service.AppParam: []string{app},
		service.EnvParam: []string{env},
	}

	res, err := get(true, false, tok, fmt.Sprintf("http://localhost:%d", port), id, params)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := id, res.Id; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}

	if want, got := app, res.App; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}

	if want, got := env, res.Env; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}

	if want, got := content, res.Content; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}
}

func TestGetInsecure(t *testing.T) {
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

	tok, err := crypto.NewToken()
	if err != nil {
		t.Fatal(err)

	}

	id, app, env, content := uuid.New().String(), "dummy", "test", "notSuperS3cret"

	crypter := &pgp.Crypter{Token: []byte(tok)}
	cypher, err := crypter.Encrypt([]byte(content))
	if err != nil {
		t.Fatal(err)
	}

	sample := fmt.Sprintf(`{"id":"%s","app_name":"%s","env":"%s","content":"%s"}`, id, app, env, string(cypher))
	src, err := models.ParseSecret(sample)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UnixNano()
	rec := &models.Record{
		Secret:    src,
		Created:   now,
		CreatedBy: "tester",
		Updated:   now,
		UpdatedBy: "tester",
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
		service.AppParam: []string{app},
		service.EnvParam: []string{env},
	}

	if _, err := get(true, false, tok, fmt.Sprintf("https://localhost:%d", port), id, params); err != nil && !strings.HasSuffix(err.Error(), "x509: certificate signed by unknown authority") {
		t.Fatal(err)
	}

	res, err := get(true, true, tok, fmt.Sprintf("https://localhost:%d", port), id, params)
	if err != nil {
		t.Fatal(err)
	}

	if want, got := id, res.Id; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}

	if want, got := app, res.App; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}

	if want, got := env, res.Env; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}

	if want, got := content, res.Content; want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}
}
