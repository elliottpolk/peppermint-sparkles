// Copyright 2016 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package datastore

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/boltdb/bolt"
)

const (
	DSFile string = "/var/lib/confgr/confgr.db"
	Bucket string = "configs"
)

var ds *bolt.DB

//  Start will open a bolt key / value. If the datastore does not currently
//  exist, a new one will be generated. It will also create the relevant bucket
//  for storage if it does not exist.
func Start() error {
	var err error

	ds, err = bolt.Open(DSFile, 0600, nil)
	if err != nil {
		return err
	}

	//  capture a shutdown and close the datastore
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL
		sig := <-c
		fmt.Printf("caught signal %s: shutting down.", sig)

		if err := ds.Close(); err != nil {
			fmt.Printf("error while closing datastore: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("datastore closed")
		os.Exit(0)
	}(sigc)

	return ds.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(Bucket)); err != nil {
			return err
		}
		return nil
	})
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
