// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/backend/datastore.go
//
package backend

import (
	"crypto/sha256"
	"fmt"
)

const (
	Redis string = "redis"
	File  string = "file"
)

type Value map[string]string

type Datastore interface {
	Close() error
	Keys() []string
	List() ([]Value, error)
	Set(key, value string) error
	Get(key string) string
	Remove(key string) error
	AddHistory(value string) error
	Historical() ([]Value, error)
}

func Key(app, env string, values ...string) string {
	in := make([]byte, 0)
	for _, v := range values {
		in = append(in, []byte(v)...)
	}

	return fmt.Sprintf("%x%s", sha256.Sum256(in), KeySuffix(app, env))
}

func KeySuffix(app, env string) string {
	suf := make([]byte, len(app)+len(env))
	copy(suf, []byte(app))
	copy(suf, []byte(env))

	return fmt.Sprintf("%x", sha256.Sum256(suf))
}
