package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"

	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
)

func TestPost(t *testing.T) {
	port := freeport()

	sample := fmt.Sprintf(`{"id":"%s","app_name":"dummy","env":"test","content":"notSuperS3cret"}`, uuid.New().String())
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

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = Handle(mux, &Handler{Backend: ds})

		wg.Done()
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
	}(ds)

	wg.Wait()

	res, err := http.Post(fmt.Sprintf("http://localhost:%d%s?%s=tester", port, PathSecrets, UserParam), "application/json", strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if code, msg := res.StatusCode, string(b); code != http.StatusCreated {
		t.Fatalf("test service POST responded with status code %d and message %s", code, msg)
	}

	if want, got := strings.TrimSpace(sample), strings.TrimSpace(string(b)); want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}
}

func TestInvalidIdPost(t *testing.T) {
	port := freeport()

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

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = Handle(mux, &Handler{Backend: ds})

		wg.Done()
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
	}(ds)

	wg.Wait()

	type sample struct {
		name    string
		value   string
		code    int
		message string
	}

	samples := []*sample{
		&sample{
			name:    "invalid_id",
			value:   `{"app_name":"dummy","env":"test","content":"notSuperS3cret"}`,
			code:    http.StatusBadRequest,
			message: "an ID for the secret must be specified",
		},
		&sample{
			name:    "invalid_app",
			value:   fmt.Sprintf(`{"id":"%s","env":"test","content":"notSuperS3cret"}`, uuid.New().String()),
			code:    http.StatusBadRequest,
			message: "an app name for the secret must be specified",
		},
		&sample{
			name:    "invalid_env",
			value:   fmt.Sprintf(`{"id":"%s","app_name":"dummy","content":"notSuperS3cret"}`, uuid.New().String()),
			code:    http.StatusBadRequest,
			message: "an environment for the secret must be specified",
		},
	}

	for _, s := range samples {
		res, err := http.Post(fmt.Sprintf("http://localhost:%d%s?%s=tester", port, PathSecrets, UserParam), "application/json", strings.NewReader(s.value))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if code, msg := res.StatusCode, strings.TrimSpace(string(b)); code != s.code || msg != s.message {
			t.Errorf("test service POST responded with status code %d and message %s for test item %s", code, msg, s.name)
		}
	}
}
