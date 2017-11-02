// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elliottpolk/confgr/datastore"

	"github.com/pkg/errors"
)

type Config struct {
	App         string `json:"app"`
	Environment string `json:"environment"`
	Value       string `json:"config"`
}

const DefaultEnv = "default"

func NewConfig(raw string) (*Config, error) {
	cfg := &Config{}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal raw config")
	}

	return cfg, nil
}

func Find(app, env string) []*Config {
	if len(app) > 0 && len(env) > 0 {
		return []*Config{
			{
				App:         app,
				Environment: env,
				Value:       datastore.Get(datastore.Key(app, env)),
			},
		}
	}

	cfgs := make([]*Config, 0)
	for _, key := range datastore.GetKeys() {
		el := strings.Split(key, "_")

		//	if the app is specified and key does not contain the app, skip
		if len(app) > 0 && app != el[0] {
			continue
		}

		cfgs = append(cfgs, &Config{
			App:         el[0],
			Environment: el[1],
			Value:       datastore.Get(key),
		})
	}

	return cfgs
}

//  Save adds the config value to the datastore using the app name as the key
func (c *Config) Save() error {
	if len(c.App) < 1 {
		return errors.New("must speicify a valid app name")
	}
	if len(c.Environment) < 1 {
		c.Environment = DefaultEnv
	}

	return datastore.Set(datastore.Key(c.App, c.Environment), c.Value)
}

//  Remove attempts to delete relevant config for the provided app name. If no
//  config exists, no error is returned.
func Remove(app, env string) error {
	if len(app) < 0 {
		return errors.New("must speicify a valid app name")
	}

	if len(env) > 0 {
		return datastore.Remove(datastore.Key(app, env))
	}

	//	need to loop through and delete all containing the app name
	for _, key := range datastore.GetKeys() {
		if app == strings.Split(key, "_")[0] {
			if err := datastore.Remove(key); err != nil {
				return errors.Wrapf(err, "unable to remove config for %s - %s", app, strings.Split(key, "_")[1])
			}
		}
	}

	return nil
}

func (cfg *Config) String() (string, error) {
	out, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (cfg *Config) MustString() string {
	str, err := cfg.String()
	if err != nil {
		return fmt.Sprintf("%+v", cfg)
	}

	return str
}
