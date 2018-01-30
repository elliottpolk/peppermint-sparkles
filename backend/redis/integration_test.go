//  +build integration

// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/backend/redis/integration_test.go
//
package redis

import (
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func getPort() string {
	var min, max int64 = 2000, 9999
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Int63n(max-min)+min)
}

func boot(name, port string) error {
	n := fmt.Sprintf("--name=%s", name)
	p := fmt.Sprintf("-p=%s:6379", port)

	return exec.Command("docker", "run", "-d", n, p, "redis:alpine").Run()
}

func kill(name string) error {
	if err := exec.Command("docker", "stop", name).Run(); err != nil {
		return err
	}
	return exec.Command("docker", "rm", name).Run()
}

func TestOpen(t *testing.T) {
	name := fmt.Sprintf("redis_%d", time.Now().UnixNano())
	port := getPort()
	if err := boot(name, port); err != nil {
		t.Fatal(err)
	}
	defer kill(name)

	ds, err := Open(&redis.Options{Addr: fmt.Sprintf("localhost:%s", port)})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()
}

func TestClose(t *testing.T) {
	name := fmt.Sprintf("redis_%d", time.Now().UnixNano())
	port := getPort()
	if err := boot(name, port); err != nil {
		t.Fatal(err)
	}
	defer kill(name)

	ds, err := Open(&redis.Options{Addr: fmt.Sprintf("localhost:%s", port)})
	if err != nil {
		t.Fatal(err)
	}

	if err := ds.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestKeys(t *testing.T) {
	name := fmt.Sprintf("redis_%d", time.Now().UnixNano())
	port := getPort()
	if err := boot(name, port); err != nil {
		t.Fatal(err)
	}
	defer kill(name)

	ds, err := Open(&redis.Options{Addr: fmt.Sprintf("localhost:%s", port)})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()

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
	name := fmt.Sprintf("redis_%d", time.Now().UnixNano())
	port := getPort()
	if err := boot(name, port); err != nil {
		t.Fatal(err)
	}
	defer kill(name)

	ds, err := Open(&redis.Options{Addr: fmt.Sprintf("localhost:%s", port)})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()

	key, bar := "foo", "bar"
	if got := ds.Get(key); len(got) > 0 {
		t.Error("expected an empty result but returned", got)
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
	name := fmt.Sprintf("redis_%d", time.Now().UnixNano())
	port := getPort()
	if err := boot(name, port); err != nil {
		t.Fatal(err)
	}
	defer kill(name)

	ds, err := Open(&redis.Options{Addr: fmt.Sprintf("localhost:%s", port)})
	if err != nil {
		t.Fatal(err)
	}
	defer ds.Close()

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
	t.Error("TODO...")
}
