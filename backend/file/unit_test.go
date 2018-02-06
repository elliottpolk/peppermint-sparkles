// Created by Elliott Polk on 24/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/backend/file/unit_test.go
//

package file

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	bolt "github.com/coreos/bbolt"
)

func TestOpen(t *testing.T) {
	what := fmt.Sprintf("psparkles_testing_%d.db", time.Now().UnixNano())
	ds, err := Open(what, &bolt.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()
	defer os.RemoveAll(what)

	key, want := "foo", "bar"
	if err := ds.Set(key, want); err != nil {
		t.Fatal(err)
	}

	if got := ds.Get(key); want != got {
		t.Errorf("\nwant %s\ngot %s\n", want, got)
	}
}

func TestClose(t *testing.T) {
	what := fmt.Sprintf("psparkles_testing_%d.db", time.Now().UnixNano())
	ds, err := Open(what, &bolt.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(what)

	key := "foo"
	if err := ds.Set(key, "bar"); err != nil {
		t.Fatal(err)
	}

	//	close to later test if data will be retrieved
	if err := ds.Close(); err != nil {
		t.Fatal(err)
	}

	if val := ds.Get(key); len(val) > 0 {
		t.Error("datastore returned value after closure")
	}
}

func TestKeys(t *testing.T) {
	what := fmt.Sprintf("psparkles_testing_%d.db", time.Now().UnixNano())
	ds, err := Open(what, &bolt.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()
	defer os.RemoveAll(what)

	if keys := ds.Keys(); len(keys) > 0 {
		t.Error("keys should have been empty")

		//	empty, just in case
		for _, k := range keys {
			if err := ds.Remove(k); err != nil {
				t.Fatal(err)
			}
		}
	}

	wants := []string{"foo", "bar", "baz"}
	for _, w := range wants {
		r := rand.NewSource(time.Now().UnixNano()).Int63()
		if err := ds.Set(w, fmt.Sprintf("%d", r)); err != nil {
			t.Fatal(err)
		}
	}

	keys := ds.Keys()
	if got, want := len(keys), len(wants); want != got {
		t.Errorf("\nwant %d\ngot %d\n", want, got)
	}

	for _, want := range wants {
		found := false
		for _, got := range keys {
			if want == got {
				found = true
			}
		}

		if !found {
			t.Error("missing key", want)
		}
	}
}

func TestSetGet(t *testing.T) {
	what := fmt.Sprintf("psparkles_testing_%d.db", time.Now().UnixNano())
	ds, err := Open(what, &bolt.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()
	defer os.RemoveAll(what)

	key, bar := "foo", "bar"
	if got := ds.Get(key); len(got) > 0 {
		t.Error("expected an empty result but return", got)
	}

	if err := ds.Set(key, bar); err != nil {
		t.Fatal(err)
	}

	if got := ds.Get(key); bar != got {
		t.Errorf("\nwant %s\ngot %s\n", bar, got)
	}

	baz := "baz"
	if err := ds.Set(key, baz); err != nil {
		t.Fatal(err)
	}

	got := ds.Get(key)
	switch got {
	case bar:
		t.Error("value remains the same after an update")

	case baz:
		break

	default:
		t.Error("Get returned unknown value of", got)
	}
}

func TestRemove(t *testing.T) {
	what := fmt.Sprintf("psparkles_testing_%d.db", time.Now().UnixNano())
	ds, err := Open(what, &bolt.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()
	defer os.RemoveAll(what)

	key, bar := "foo", "bar"
	if err := ds.Set(key, bar); err != nil {
		t.Fatal(err)
	}

	//	test if set worked
	if got := ds.Get(key); bar != got {
		t.Errorf("\nwant %s\ngot %s\n", bar, got)
	}

	if err := ds.Remove(key); err != nil {
		t.Fatal(err)
	}

	if got := ds.Get(key); len(got) > 0 || got == bar {
		t.Errorf("non-empty value of %s was returned after removal", got)
	}
}

func TestList(t *testing.T) {
	what := fmt.Sprintf("psparkles_testing_%d.db", time.Now().UnixNano())
	ds, err := Open(what, &bolt.Options{})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()
	defer os.RemoveAll(what)

	wants := map[string]string{
		"97df3588b5a3f24babc3851b372f0ba71a9dcdded43b14b9d06961bfc1707d9d": "foo",
		"ebf3b019bb7e36bdc0fbc4159345c04af54193ccf43ae4572922f6d4aa94bd5b": "bar",
		"2c60dbf3773104dce76dfbda9b82a729e98a42a7a0b3f9bae5095c7bed752b90": "baz",
		"796362b8b4289fca4d666ab486487d6699e828f9c098fc1c91566c291ef682f6": "biz",
	}

	//	fill the datastore with sample values
	for k, v := range wants {
		if err := ds.Set(k, v); err != nil {
			t.Fatal(err)
		}
	}

	gots, err := ds.List()
	if err != nil {
		t.Fatal(err)
	}

	for _, vals := range gots {
		for k, got := range vals {
			if want, ok := wants[k]; !ok || want != got {
				t.Errorf("\nwant %s\ngot %s\n---\n", want, got)
			}
		}
	}
}
