// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/peppermint-sparkles/backend/file/datastore.go
//
package file

import (
	"git.platform.manulife.io/oa-montreal/peppermint-sparkles/backend"

	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
)

const bucket string = "mints"

type Datastore struct {
	db *bolt.DB
}

var ErrInvalidDatastore error = errors.New("no valid datastore")

func Open(name string, opts *bolt.Options) (*Datastore, error) {
	db, err := bolt.Open(name, 0600, opts)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open datastore file")
	}

	ds := &Datastore{db: db}

	//  ensure that the bucket exists
	err = ds.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to ensure creation of bucket")
	}

	return ds, nil
}

//  Close attempts to close the file of the datastore
func (ds *Datastore) Close() error {
	if ds.db != nil {
		return ds.db.Close()
	}
	return nil
}

//  Keys iterates over the available keys and returns as a list.
func (ds *Datastore) Keys() []string {
	keys := make([]string, 0)
	if ds.db != nil {
		ds.db.View(func(tx *bolt.Tx) error {
			curs := tx.Bucket([]byte(bucket)).Cursor()
			for k, _ := curs.First(); k != nil; k, _ = curs.Next() {
				keys = append(keys, string(k))
			}

			return nil
		})
	}

	return keys
}

//  Set adds a new entry into the key/value store. If the key exists, the old
//  value will be overwritten.
func (ds *Datastore) Set(key, value string) error {
	if ds.db == nil {
		return ErrInvalidDatastore
	}
	return ds.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Put([]byte(key), []byte(value))
	})
}

//  Get retrieves the relevant content for the provided key.
func (ds *Datastore) Get(key string) string {
	if ds.db == nil {
		return ""
	}

	var val string
	ds.db.View(func(tx *bolt.Tx) error {
		val = string(tx.Bucket([]byte(bucket)).Get([]byte(key)))
		return nil
	})

	return val
}

//  Remove deletes the content for the provided key. No error is returned if the
//  provided key does not exist.
func (ds *Datastore) Remove(key string) error {
	if ds.db == nil {
		return ErrInvalidDatastore
	}

	return ds.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(bucket)).Delete([]byte(key))
	})
}

func (ds *Datastore) List() ([]backend.Value, error) {
	if ds.db == nil {
		return nil, ErrInvalidDatastore
	}

	vals := make([]backend.Value, 0)
	for _, k := range ds.Keys() {
		vals = append(vals, backend.Value{k: ds.Get(k)})
	}

	return vals, nil
}
