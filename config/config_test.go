// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import (
	"encoding/json"
	"testing"

	"github.com/elliottpolk/confgr/datastore"
)

const (
	want string = `{"app":"foo","environment":"test","config":"{\"user\":\"test\",\"password\":\"Sup3rS3(r37\"}"}`

	app  string = "foo"
	env  string = "test"
	conf string = `{"user":"test","password":"Sup3rS3(r37"}`

	dbfile string = "confgr_test.db"
)

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfig(want)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := json.Marshal(cfg)
	if err != nil {
		t.Error(err)
		return
	}

	if got := string(out); want != got {
		t.Error("\nwanted %s\ngot %s\n", want, got)
	}
}

func TestFind(t *testing.T) {
	ds, err := datastore.Open(dbfile)
	if err != nil {
		t.Error(err)
		return
	}
	defer ds.Close(true)

	cfg, err := NewConfig(want)
	if err != nil {
		t.Error(err)
		return
	}

	if err := cfg.Save(); err != nil {
		t.Error(err)
		return
	}

	res := Find("", "")
	if cnt := len(res); cnt != 1 {
		t.Errorf("results for find should be exactly 1, found %d", cnt)
		return
	}

	c := res[0]
	if want, got := app, c.App; want != got {
		t.Errorf("wanted %s - got %s", want, got)
	}

	if want, got := env, c.Environment; want != got {
		t.Errorf("wanted %s - got %s", want, got)
	}

	if want, got := conf, c.Value; want != got {
		t.Errorf("wanted %s - got %s", want, got)
	}
}
