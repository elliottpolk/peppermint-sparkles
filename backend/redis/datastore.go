// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/backend/redis/datastore.go
//
package redis

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"git.platform.manulife.io/go-common/log"
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var ErrInvalidDatastore error = errors.New("no valid datastore")

type Datastore struct {
	client     *redis.Client
	historical *redis.Client
}

func Open(opts *redis.Options) (*Datastore, error) {
	opts.DB = 0 //	ensure that is uses the same (default) DB every time
	ds := &Datastore{client: redis.NewClient(opts)}

	opts.DB = 1 //	ensure that is uses the same (default) historical DB every time
	ds.historical = redis.NewClient(opts)

	//	ensure a valid connection prior to returning
	if _, err := ds.client.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "unable to ping redis datastore")
	}

	//	ensure a valid connection to the historical store prior to returning
	if _, err := ds.historical.Ping().Result(); err != nil {
		return nil, errors.Wrap(err, "unable to ping historical redis datastore")
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

func keys(client *redis.Client) ([]string, error) {
	keys := make([]string, 0)
	if client != nil {
		k, err := client.Keys("*").Result()
		if err != nil {
			return keys, err
		}

		for _, v := range k {
			keys = append(keys, v)
		}
	}
	return keys, nil
}

func (ds *Datastore) Keys() []string {
	vals, err := keys(ds.client)
	if err != nil {
		//	log error but still return the empty list
		log.Error(err, "unable to retrieve keys")
	}
	return vals
}

func (ds *Datastore) Set(key, value string) error {
	if ds.client == nil {
		return ErrInvalidDatastore
	}
	return ds.client.Set(key, value, 0).Err()
}

func get(key string, client *redis.Client) string {
	res, err := client.Get(key).Result()
	if err != nil && err != redis.Nil {
		//	log error but still return the empty string
		log.Error(err, "unable to retrieve result for key")
		return ""
	}
	return res
}

func (ds *Datastore) Get(key string) string {
	return get(key, ds.client)
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

func (ds *Datastore) AddHistory(value string) error {
	if ds.historical == nil {
		return ErrInvalidDatastore
	}

	buf := make([]byte, 2048)
	if _, err := rand.Read(buf); err != nil {
		return errors.Wrap(err, "unable to read in random data to generate key")
	}

	//	generate SHA256 token from random content to be stored + random data to
	// 	attempt to prevent collisions
	key := fmt.Sprintf("%x", sha256.Sum256(append([]byte(value), buf...)))

	return ds.historical.Set(key, value, 0).Err()
}

func (ds *Datastore) historicalKeys() []string {
	vals, err := keys(ds.historical)
	if err != nil {
		//	log error but still return the empty list
		log.Error(err, "unable to retrieve historical keys")
	}
	return vals
}

func (ds *Datastore) Historical() ([]backend.Value, error) {
	if ds.historical == nil {
		return nil, ErrInvalidDatastore
	}

	vals := make([]backend.Value, 0)
	for _, k := range ds.historicalKeys() {
		if res := get(k, ds.historical); len(res) > 0 {
			vals = append(vals, backend.Value{k: res})
		}
	}

	return vals, nil
}
