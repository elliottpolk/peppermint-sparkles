package models

import (
	"encoding/json"

	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"

	"github.com/pkg/errors"
)

const (
	ActiveStatus  string = "active"
	ArchiveStatus string = "archived"
	InvalidStatus string = "invalid"
)

type Record struct {
	*Secret `json:"secret"`

	Created   int64  `json:"created"`
	CreatedBy string `json:"created_by"`
	Updated   int64  `json:"updated"`
	UpdatedBy string `json:"updated_by"`
	Status    string `json:"status"`
}

func ParseRecord(raw string) (*Record, error) {
	r := &Record{}
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return nil, errors.Wrap(err, "unable to parse raw record")
	}
	return r, nil
}

func (r *Record) Exists(in backend.Datastore) bool {
	return len(in.Get(r.Secret.Id)) > 0
}

func (r *Record) Write(where backend.Datastore) error {
	out, err := r.String()
	if err != nil {
		return errors.Wrap(err, "unable to prep record for storage")
	}
	return where.Set(r.Secret.Id, out)
}

func (r *Record) String() (string, error) {
	out, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (r *Record) MustString() string {
	str, err := r.String()
	if err != nil {
		return ""
	}

	return str
}
