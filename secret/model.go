// Created by Elliott Polk on 25/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/secret/model.go
//
package secret

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Secret struct {
	Id      string `json:"id,omitempty"`
	App     string `json:"app_name"`
	Env     string `json:"env"`
	Content string `json:"content"`
}

func NewSecret(raw string) (*Secret, error) {
	s := &Secret{}

	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal raw secret")
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
		return fmt.Sprintf("%+v", s)
	}

	return str
}
