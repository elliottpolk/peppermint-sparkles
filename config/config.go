// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import "github.com/elliottpolk/confgr/datastore"

type Config struct {
	App   string `json:"app"`
	Value string `json:"config"`
}

//  ListApps retrieves the list of keys from the datastore
func ListApps() []string {
	return datastore.GetKeys()
}

//  Save adds the config value to the datastore using the app name as the key
func (c *Config) Save() error {
	return datastore.Set(c.App, c.Value)
}

//  Get retrieves the config for the provided app name. If no config exists, an
//  empty string is set for the Config.Value.
func Get(app string) *Config {
	return &Config{app, datastore.Get(app)}
}

//  Remove attempts to delete relevant config for the provided app name. If no
//  config exists, no error is returned.
func Remove(app string) error {
	return datastore.Remove(app)
}
