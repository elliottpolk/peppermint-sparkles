package file

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"github.com/manulife-gwam/peppermint-sparkles/backend"

	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
)

const (
	bucket     string = "mints"
	historical string = "buttermints"
)

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

	//  ensure that the buckets exist
	for _, b := range []string{bucket, historical} {
		err = ds.db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte(b)); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to ensure creation of bucket %s", b)
		}
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

func (ds *Datastore) keys(b string) []string {
	vals := make([]string, 0)
	if ds.db != nil {
		ds.db.View(func(tx *bolt.Tx) error {
			curs := tx.Bucket([]byte(b)).Cursor()
			for k, _ := curs.First(); k != nil; k, _ = curs.Next() {
				vals = append(vals, string(k))
			}

			return nil
		})
	}
	return vals
}

//  Keys iterates over the available keys and returns as a list.
func (ds *Datastore) Keys() []string {
	return ds.keys(bucket)
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

func (ds *Datastore) get(b, k string) string {
	if ds.db == nil {
		return ""
	}
	var val string
	ds.db.View(func(tx *bolt.Tx) error {
		val = string(tx.Bucket([]byte(b)).Get([]byte(k)))
		return nil
	})

	return val
}

//  Get retrieves the relevant content for the provided key.
func (ds *Datastore) Get(key string) string {
	return ds.get(bucket, key)
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

func (ds *Datastore) AddHistory(value string) error {
	if ds.db == nil {
		return ErrInvalidDatastore
	}

	buf := make([]byte, 2048)
	if _, err := rand.Read(buf); err != nil {
		return errors.Wrap(err, "unable to read in random data to generate key")
	}

	//	generate SHA256 token from random content to be stored + random data to
	// 	attempt to prevent collisions
	key := fmt.Sprintf("%x", sha256.Sum256(append([]byte(value), buf...)))

	return ds.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(historical)).Put([]byte(key), []byte(value))
	})
}

func (ds *Datastore) historicalKeys() []string {
	return ds.keys(historical)
}

func (ds *Datastore) Historical() ([]backend.Value, error) {
	if ds.db == nil {
		return nil, ErrInvalidDatastore
	}

	vals := make([]backend.Value, 0)
	for _, k := range ds.historicalKeys() {
		if res := ds.get(historical, historical); len(res) > 0 {
			vals = append(vals, backend.Value{k: res})
		}
	}

	return vals, nil
}
