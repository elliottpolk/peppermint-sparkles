package models

import (
	"encoding/json"

	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"

	"github.com/pkg/errors"
)

const (
	CreateAction  string = "create"
	UpdateAction  string = "update"
	ArchiveAction string = "archive"
)

type Historical struct {
	*Record `json:"record"`

	Action    string `json:"action"`
	Created   int64  `json:"created"`
	CreatedBy string `json:"created_by"`
}

func FromCurrent(what string) (*Historical, error) {
	r, err := ParseRecord(what)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse current record")
	}

	return &Historical{Record: r}, nil
}

func (h *Historical) Write(where backend.Datastore, why, who string, when int64) error {
	h.Action = why
	h.CreatedBy = who
	h.Created = when

	out, err := h.String()
	if err != nil {
		return errors.Wrap(err, "unable to prep historical for storage")
	}

	return where.AddHistory(out)
}

func (h *Historical) String() (string, error) {
	out, err := json.MarshalIndent(h, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (h *Historical) MustString() string {
	str, err := h.String()
	if err != nil {
		return ""
	}

	return str
}
