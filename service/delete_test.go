package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	fileds "github.com/manulife-gwam/peppermint-sparkles/backend/file"
	"github.com/manulife-gwam/peppermint-sparkles/models"

	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
)

func TestDelete(t *testing.T) {
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

	app, env, usr := "dummy", "test", "tester"

	sample := fmt.Sprintf(`{"id":"%s","app_name":"%s","env":"%s","content":"notSuperS3cret"}`, uuid.New().String(), app, env)
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

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:%d%s/%s", port, PathSecrets, src.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	params := &url.Values{
		AppParam:  []string{src.App},
		EnvParam:  []string{src.Env},
		UserParam: []string{usr},
	}
	req.URL.RawQuery = params.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if code, msg := res.StatusCode, string(b); code != http.StatusOK {
		t.Fatalf("test service DELETE responded with status code %d and message %s", code, msg)
	}

	if raw := ds.Get(src.Id); len(raw) > 0 {
		t.Errorf("the deleted secert id responded with %s", raw)
	}
}
