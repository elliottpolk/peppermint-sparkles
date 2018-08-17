package service

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

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

func TestApi(t *testing.T) {
	const sample string = `{
 "id": "6f0f9805-08c6-48f2-b3c4-fe8e7c35ea4a",
 "app_name": "dummy",
 "env": "test",
 "content": "notSuperS3cret"
}`

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

	// set a wait group to allow for some setup time
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ds *fileds.Datastore) {
		mux := http.NewServeMux()
		mux = Handle(mux, &Handler{Backend: ds})

		wg.Done()
		t.Fatal(http.ListenAndServe(":6000", mux))
	}(ds)

	wg.Wait()

	// test POST
	r0, err := http.Post(fmt.Sprintf("http://localhost:6000/api/v2/secrets?%s=tester", UserParam), "application/json", strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	defer r0.Body.Close()

	pb, err := ioutil.ReadAll(r0.Body)
	if err != nil {
		t.Fatal(err)
	}

	if code, msg := r0.StatusCode, string(pb); code != http.StatusCreated {
		t.Fatalf("test service POST responded with status code %d and message %s", code, msg)
	}

	// test GET
	r1, err := http.Get(fmt.Sprintf("http://localhost:6000/api/v2/secrets/%s", src.Id))
	if err != nil {
		t.Fatal(err)
	}
	defer r1.Body.Close()

	gb, err := ioutil.ReadAll(r1.Body)
	if err != nil {
		t.Fatal(err)
	}

	if code, msg := r1.StatusCode, string(gb); code != http.StatusOK {
		t.Fatalf("test service GET responded with status code %d and message %s", code, msg)
	}

	s, err := models.ParseSecret(string(gb))
	if err != nil {
		t.Fatal(err)
	}

	if want, got := sample, s.MustString(); want != got {
		t.Errorf("want %s\n\ngot %s\n", want, got)
	}
}
