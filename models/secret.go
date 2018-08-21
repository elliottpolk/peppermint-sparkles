package models

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const tag string = "peppermint-sparkles.models"

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
