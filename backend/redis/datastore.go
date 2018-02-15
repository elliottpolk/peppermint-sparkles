// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/backend/redis/datastore.go
//
package redis

import (
	"gitlab.manulife.com/go-common/log"
	"gitlab.manulife.com/oa-montreal/peppermint-sparkles/backend"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var ErrInvalidDatastore error = errors.New("no valid datastore")

type Datastore struct {
	client *redis.Client
}

func Open(opts *redis.Options) (*Datastore, error) {
	opts.DB = 0 //	ensure that is uses the same (default) DB every time
	ds := &Datastore{client: redis.NewClient(opts)}

	//	ensure a valid connection prior to returning the Datastore
	if _, err := ds.client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "unable to ping redis datastore")
	}

	return ds, nil
}

//  Close attempts to close the redis connection of the datastore
func (ds *Datastore) Close() error {
	if ds.client != nil {
		return ds.client.Close()
	}
	return nil
}

func (ds *Datastore) Keys() []string {
	keys := make([]string, 0)
	if ds.client != nil {
		k, err := ds.client.Keys("*").Result()
		if err != nil {
			//	log error but still return the empty list
			log.Error(err, "unable to retrieve keys")
			return keys
		}

		for _, v := range k {
			keys = append(keys, v)
		}
	}

	return keys
}

func (ds *Datastore) Set(key, value string) error {
	if ds.client == nil {
		return ErrInvalidDatastore
	}
	return ds.client.Set(key, value, 0).Err()
}

func (ds *Datastore) Get(key string) string {
	if ds.client != nil {
		res, err := ds.client.Get(key).Result()
		if err != nil && err != redis.Nil {
			//	log error but still return the empty string
			log.Error(err, "unable to retrieve result for key")
			return ""
		}
		return res
	}

	return ""
}

func (ds *Datastore) Remove(key string) error {
	if ds.client == nil {
		return ErrInvalidDatastore
	}

	return ds.client.Del(key).Err()
}

func (ds *Datastore) List() ([]backend.Value, error) {
	if ds.client == nil {
		return nil, ErrInvalidDatastore
	}

	vals := make([]backend.Value, 0)
	for _, k := range ds.Keys() {
		vals = append(vals, backend.Value{k: ds.Get(k)})
	}

	return vals, nil
}
