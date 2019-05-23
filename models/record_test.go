package models

import (
	"fmt"
	"os"
	"testing"

	fileds "git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend/file"

	bolt "github.com/coreos/bbolt"
	"github.com/google/uuid"
)

func TestRecord(t *testing.T) {
	const sample string = `{
 "secret": {
  "id": "6f0f9805-08c6-48f2-b3c4-fe8e7c35ea4a",
  "app_name": "dummy",
  "env": "test",
  "content": "notSuperS3cret"
 },
 "created": 1534474065732344471,
 "created_by": "tester",
 "updated": 1534474065732344471,
 "updated_by": "tester",
 "status": "active"
}`

	tmpRepo := fmt.Sprintf("test_%s.db", uuid.New().String())

	ds, err := fileds.Open(tmpRepo, bolt.DefaultOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer func(ds *fileds.Datastore) {
		ds.Close()
		if err := os.RemoveAll(tmpRepo); err != nil {
			t.Errorf("unable to remove temporary test repo %s\n", tmpRepo)
		}
	}(ds)

	// tests parsing of record
	r, err := ParseRecord(sample)
	if err != nil {
		t.Fatal(err)
	}

	// tests stringer of record
	if want, got := sample, r.MustString(); want != got {
		t.Errorf("want: %s\n\ngot: %s", want, got)
	}

	// tests writing of record to datastore
	if err := r.Write(ds); err != nil {
		t.Fatal(err)
	}

	// tests a valid records existence
	if !r.Exists(ds) {
		t.Error("failed to verify test item written")
	}

	// tests a valid records existence
	invalid := &Record{Secret: &Secret{Id: "invalid_not_real_id"}}
	if invalid.Exists(ds) {
		t.Error("found invalid record")
	}
}
