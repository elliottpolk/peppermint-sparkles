package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/crypto/pgp"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/service"

	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
)

func TestSet(t *testing.T) {
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

	tok, err := crypto.NewToken()
	if err != nil {
		t.Fatal(err)
	}

	app, env, content := "dummy", "test", "notSuperS3cret"
	raw := fmt.Sprintf(`{"app_name":"%s","env":"%s","content":"%s"}`, app, env, content)
	addr := fmt.Sprintf("http://localhost:%d", port)

	s, err := set(true, false, tok, "tester", raw, addr)
	if err != nil {
		t.Fatal(err)
	}

	//	decrypt content with provided token
	txt, err := (&pgp.Crypter{Token: []byte(tok)}).Decrypt([]byte(s.Content))
	if err != nil {
		t.Fatal(err)
	}

	if want, got := content, string(txt); want != got {
		t.Errorf("want %s\ngot  %s", want, got)
	}

	// verify app name was returned properly
	if want, got := app, s.App; want != got {
		t.Errorf("want %s\ngot  %s", want, got)
	}

	// verify environment was returned properly
	if want, got := env, s.Env; want != got {
		t.Errorf("want %s\ngot  %s", want, got)
	}
}

func TestSetInsecure(t *testing.T) {
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

	tok, err := crypto.NewToken()
	if err != nil {
		t.Fatal(err)
	}

	app, env, content := "dummy", "test", "notSuperS3cret"
	raw := fmt.Sprintf(`{"app_name":"%s","env":"%s","content":"%s"}`, app, env, content)
	addr := fmt.Sprintf("https://localhost:%d", port)

	if _, err := set(true, false, tok, "tester", raw, addr); err != nil && !strings.HasSuffix(err.Error(), "x509: certificate signed by unknown authority") {
		t.Fatal(err)
	}

	s, err := set(true, true, tok, "tester", raw, addr)
	if err != nil {
		t.Fatal(err)
	}

	//	decrypt content with provided token
	txt, err := (&pgp.Crypter{Token: []byte(tok)}).Decrypt([]byte(s.Content))
	if err != nil {
		t.Fatal(err)
	}

	if want, got := content, string(txt); want != got {
		t.Errorf("want %s\ngot  %s", want, got)
	}

	// verify app name was returned properly
	if want, got := app, s.App; want != got {
		t.Errorf("want %s\ngot  %s", want, got)
	}

	// verify environment was returned properly
	if want, got := env, s.Env; want != got {
		t.Errorf("want %s\ngot  %s", want, got)
	}
}

func TestInvalidSet(t *testing.T) {
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

	tok, err := crypto.NewToken()
	if err != nil {
		t.Fatal(err)
	}

	app, env, content := "dummy", "test", "notSuperS3cret"
	addr := fmt.Sprintf("http://localhost:%d", port)

	type sample struct {
		name    string
		value   string
		message string
	}

	samples := []*sample{
		&sample{
			name:    "invalid_app",
			value:   fmt.Sprintf(`{"env":"%s","content":"%s"}`, env, content),
			message: "unable to send secret: secrets service responded with status code 400 and message an app name for the secret must be specified",
		},
		&sample{
			name:    "invalid_env",
			value:   fmt.Sprintf(`{"app_name":"%s","content":"%s"}`, app, content),
			message: "unable to send secret: secrets service responded with status code 400 and message an environment for the secret must be specified",
		},
	}

	for _, s := range samples {
		if _, err := set(true, false, tok, "tester", s.value, addr); err != nil && strings.TrimSpace(err.Error()) != s.message {
			t.Errorf("\nwant %s\ngot  %s\n", s.message, err.Error())
		}
	}
}
