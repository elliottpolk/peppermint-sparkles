// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import (
	"fmt"
	"strings"

	"github.com/elliottpolk/confgr/datastore"
)

type Config struct {
	App         string `json:"app"`
	Environment string `json:"environment"`
	Value       string `json:"config"`
}

const DefaultEnv = "default"

//  ListApps retrieves a map of app names and available environments from the
//	datastore. If no keys exist, an empty map is returned
func ListApps() map[string][]string {
	listing := make(map[string][]string)
	for _, k := range datastore.GetKeys() {
		name := strings.Split(k, "_")[0]
		env := strings.Split(k, "_")[1]

		item, ok := listing[name]
		if !ok {
			item = make([]string, 0)
		}
		listing[name] = append(item, env)
	}

	return listing
}

//  Save adds the config value to the datastore using the app name as the key
func (c *Config) Save() error {
	if len(c.Environment) < 1 {
		c.Environment = DefaultEnv
	}

	return datastore.Set(fmt.Sprintf("%s_%s", c.App, c.Environment), c.Value)
}

//  Get retrieves the config for the provided app name. If no config exists, an
//  empty string is set for the Config.Value
func Get(app, env string) *Config {
	if len(env) < 1 {
		env = DefaultEnv
	}

	return &Config{app, env, datastore.Get(fmt.Sprintf("%s_%s", app, env))}
}

//  Remove attempts to delete relevant config for the provided app name. If no
//  config exists, no error is returned.
func Remove(app, env string) error {
	if len(env) < 1 {
		env = DefaultEnv
	}

	return datastore.Remove(fmt.Sprintf("%s_%s", app, env))
}
