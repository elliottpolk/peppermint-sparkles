// Created by Elliott Polk on 30/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// internal/pcf/vcap/service.go
//

package vcap

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

const envServices string = "VCAP_SERVICES"

type Instance struct {
	Name           string      `json:"name"`
	Label          string      `json:"label"`
	Tags           []string    `json:"tags"`
	Plan           string      `json:"plan"`
	Credentials    Credentials `json:"credentials"`
	Provider       interface{} `json:"provider"`
	SyslogDrainURL interface{} `json:"syslog_drain_url"`
}

type Services map[string][]Instance
type Credentials map[string]interface{}

func GetServices() (Services, error) {
	env := os.Getenv(envServices)
	if len(env) < 1 {
		return nil, nil
	}

	s := Services{}
	if err := json.Unmarshal([]byte(env), &s); err != nil {
		return nil, errors.Wrap(err, "unable to parse vcap services")
	}

	return s, nil
}

func (s Services) Get(name string) []Instance {
	return s[name]
}

func (s Services) Tagged(tag string) *Instance {
	for _, list := range s {
		for _, inst := range list {
			for _, t := range inst.Tags {
				if t == tag {
					return &inst
				}
			}
		}
	}
	return nil
}
