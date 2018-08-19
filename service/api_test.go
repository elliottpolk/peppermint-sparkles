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

	"git.platform.manulife.io/go-common/log"
	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/models"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/uuid"

	bolt "github.com/coreos/bbolt"
)

func init() {
	log.InitTester()
}

func TestGetId(t *testing.T) {
	want := "6f0f9805-08c6-48f2-b3c4-fe8e7c35ea4a"
	m, got, err := getId(fmt.Sprintf("secrets/%s", want))
	if err != nil {
		t.Fatal(err)
	}

	if !m {
		t.Error("expected match and getId returned false")
	}

	if got != want {
		t.Errorf("want %s\n\ngot %s\n", want, got)
	}

	noMatch, invalidId, err := getId(fmt.Sprintf("testing/get/secrets/%s-asdfqwerty", want))
	if err != nil {
		t.Fatal(err)
	}

	if noMatch {
		t.Error("expected to not match and getId returned true")
	}

	if invalidId == want {
		t.Errorf("expected no match but returned %s", invalidId)
	}
}

func TestPost(t *testing.T) {
	const port string = "6000"

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

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = Handle(mux, &Handler{Backend: ds})

		wg.Done()
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
	}(ds)

	wg.Wait()

	res, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v2/secrets?%s=tester", port, UserParam), "application/json", strings.NewReader(sample))
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
	const port string = "6001"

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

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = Handle(mux, &Handler{Backend: ds})

		wg.Done()
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
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
			value:   fmt.Sprintf(`{"id":"%s","env":"test","content":"notSuperS3cret"}`, uuid.GetV4()),
			code:    http.StatusBadRequest,
			message: "an app name for the secret must be specified",
		},
		&sample{
			name:    "invalid_env",
			value:   fmt.Sprintf(`{"id":"%s","app_name":"dummy","content":"notSuperS3cret"}`, uuid.GetV4()),
			code:    http.StatusBadRequest,
			message: "an environment for the secret must be specified",
		},
	}

	for _, s := range samples {
		res, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v2/secrets?%s=tester", port, UserParam), "application/json", strings.NewReader(s.value))
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

func TestGet(t *testing.T) {
	const port string = "6002"
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
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
	}(ds)

	wg.Wait()

	res, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v2/secrets/%s?%s=%s&%s=%s", port, src.Id, AppParam, src.App, EnvParam, src.Env))
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
	const port string = "6003"

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
		t.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
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
			from:    fmt.Sprintf("http://localhost:%s/api/v2/secrets/%s?%s=%s&%s=%s", port, uuid.GetV4(), AppParam, src.App, EnvParam, src.Env),
			code:    http.StatusNotFound,
			message: "file not found",
		},
		&sample{
			name:    "invalid_app_name",
			from:    fmt.Sprintf("http://localhost:%s/api/v2/secrets/%s?%s=%s&%s=%s", port, src.Id, AppParam, "flerp", EnvParam, src.Env),
			code:    http.StatusBadRequest,
			message: "app ID and name are invalid",
		},
		&sample{
			name:    "invalid_app_env",
			from:    fmt.Sprintf("http://localhost:%s/api/v2/secrets/%s?%s=%s&%s=%s", port, src.Id, AppParam, src.App, EnvParam, "PROD"),
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
