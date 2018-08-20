package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/uuid"

	bolt "github.com/coreos/bbolt"
)

func TestGet(t *testing.T) {
	port := freeport()

	sample := fmt.Sprintf(`{"id":"%s","app_name":"dummy","env":"test","content":"notSuperS3cret"}`, uuid.GetV4())
	repo := fmt.Sprintf("test_%s.db", uuid.GetV4())

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

	res, err := http.Get(fmt.Sprintf("http://localhost:%d%s/%s?%s=%s&%s=%s", port, PathSecrets, src.Id, AppParam, src.App, EnvParam, src.Env))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if code, msg := res.StatusCode, string(b); code != http.StatusOK {
		t.Fatalf("test service GET responded with status code %d and message %s", code, msg)
	}

	s, err := models.ParseSecret(string(b))
	if err != nil {
		t.Fatal(err)
	}

	if want, got := src.MustString(), s.MustString(); want != got {
		t.Errorf("\nwant %s\ngot  %s\n", want, got)
	}
}

func TestInvalidGet(t *testing.T) {
	port := freeport()

	repo := fmt.Sprintf("test_%s.db", uuid.GetV4())
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

	raw := fmt.Sprintf(`{"id":"%s","app_name":"dummy","env":"test","content":"notSuperS3cret"}`, uuid.GetV4())
	src, err := models.ParseSecret(raw)
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
		from    string
		code    int
		message string
	}

	samples := []*sample{
		&sample{
			name:    "invalid_id",
			from:    fmt.Sprintf("http://localhost:%d%s/%s?%s=%s&%s=%s", port, PathSecrets, uuid.GetV4(), AppParam, src.App, EnvParam, src.Env),
			code:    http.StatusNotFound,
			message: "file not found",
		},
		&sample{
			name:    "invalid_app_name",
			from:    fmt.Sprintf("http://localhost:%d%s/%s?%s=%s&%s=%s", port, PathSecrets, src.Id, AppParam, "flerp", EnvParam, src.Env),
			code:    http.StatusBadRequest,
			message: "app ID and name are invalid",
		},
		&sample{
			name:    "invalid_app_env",
			from:    fmt.Sprintf("http://localhost:%d%s/%s?%s=%s&%s=%s", port, PathSecrets, src.Id, AppParam, src.App, EnvParam, "PROD"),
			code:    http.StatusBadRequest,
			message: "app ID and environment are invalid",
		},
	}

	for _, s := range samples {
		res, err := http.Get(s.from)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		if code, msg := res.StatusCode, strings.TrimSpace(string(b)); code != s.code && msg != s.message {
			t.Fatalf("test service GET responded with status code %d and message %s for test item %s", code, msg, s.name)
		}
	}
}
