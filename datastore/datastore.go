// Created by Elliott Polk on 23/01/2018
// Copyright © 2018 Manulife AM. All rights reserved.
// oa-montreal/campx/main.go
//
package datastore

import (
	"fmt"
	"os"

	"git.platform.manulife.io/oa-montreal/campx/log"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

const Bucket string = "configs"

type Datastore struct {
	*bolt.DB
	file string
}

var ds *Datastore

func Current() *Datastore {
	if ds == nil {
		ds = &Datastore{}
	}
	return ds
}

func Open(f string) (*Datastore, error) {
	ds = &Datastore{file: f}
	if err := ds.Open(); err != nil {
		return nil, errors.Wrap(err, "unable to open datastore file")
	}

	return ds, nil
}

//  Start will open a bolt key / value. If the datastore does not currently
//  exist, a new one will be generated. It will also create the relevant bucket
//  for storage if it does not exist.
func (ds *Datastore) Open() error {
	db, err := bolt.Open(ds.file, 0600, nil)
	if err != nil {
		return errors.Wrap(err, "unable to open bolt data store")
	}

	ds.DB = db

	return ds.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(Bucket)); err != nil {
			return errors.Wrap(err, "unable to create bucket")
		}
		return nil
	})
}

func (ds *Datastore) Close(rm bool) error {
	if err := ds.DB.Close(); err != nil {
		return errors.Wrap(err, "unable to close datastore")
	}

	if rm {
		if err := os.RemoveAll(ds.file); err != nil {
			log.Error(err, "unable to remove datastore")
			os.Exit(1)
		}
	}

	return nil
}

func Key(values ...string) string {
	res := ""
	for _, v := range values {
		if len(res) > 0 {
			res = fmt.Sprintf("%s_", res)
		}
		res = fmt.Sprintf("%s%s", res, v)
	}

	return res
}

//  GetKeys iterates over the available keys and returns them as a list.
func GetKeys() []string {
	keys := make([]string, 0)

	ds.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(Bucket)).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}

		return nil
	})

	return keys
}

//  Set adds a new entry into the key/value store. If the key exists, the old
//  value will be overwritten.
func Set(key, value string) error {
	return ds.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(Bucket)).Put([]byte(key), []byte(value))
	})
}

//  Get retrieves the relevant content for the provided key.
func Get(key string) string {
	var res string

	ds.View(func(tx *bolt.Tx) error {
		res = string(tx.Bucket([]byte(Bucket)).Get([]byte(key)))
		return nil
	})

	return res
}

//  Remove deletes the content for the provided key. No error is returned if the
//  provided key does not exist.
func Remove(key string) error {
	return ds.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(Bucket)).Delete([]byte(key))
	})
}
