// Created by Elliott Polk on 25/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/models/secret.go
//
package models

import (
	"crypto/rand"
	"encoding/json"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"

	"github.com/pkg/errors"
)

type Secret struct {
	Id      string `json:"id,omitempty"`
	App     string `json:"app_name"`
	Env     string `json:"env"`
	Content string `json:"content"`
}

func ParseSecret(raw string) (*Secret, error) {
	s := &Secret{}

	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return nil, errors.Wrap(err, "unable to parse raw secret")
	}

	return s, nil
}

func (s *Secret) String() (string, error) {
	out, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (s *Secret) MustString() string {
	str, err := s.String()
	if err != nil {
		return ""
	}

	return str
}

func (s *Secret) NewId() string {
	buf := make([]byte, 1024)
	if _, err := rand.Read(buf); err != nil {
		log.Error(err, "unable to read in random data for id generation")
	}

	//	FIXME ... should really allow for a retry on the random read
	return backend.Key(s.App, s.Env, string(buf))
}
