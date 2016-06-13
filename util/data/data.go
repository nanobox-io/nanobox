// Package data ...
package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/boltdb/bolt"

	"github.com/nanobox-io/nanobox/util"
)

// db establishes a bolt.DB for future use
func db() *bolt.DB {

	//
	boltDB, err := bolt.Open(filepath.ToSlash(filepath.Join(util.GlobalDir(), "data.db")), 0666, nil)
	if err != nil {
		panic(err)
	}

	return boltDB
}

// Put inserts or updates an element into the bolt database
func Put(bucket, id string, v interface{}) error {

	d := db()
	defer d.Close()

	//
	return d.Update(func(tx *bolt.Tx) error {

		// Create a bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		//
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		bucket.Put([]byte(id), bytes)

		return nil
	})
}

// Get retrieves an element from the database unlike the default behavior from
// boltdb this will return an error if you try getting something that doesnt exist
func Get(bucket, id string, v interface{}) error {

	// Read value back in a different read-only transaction.
	d := db()
	defer d.Close()

	//
	return d.View(func(tx *bolt.Tx) error {

		//
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return fmt.Errorf("no record found")
		}

		//
		value := bucket.Get([]byte(id))
		if value == nil || len(value) == 0 {
			return fmt.Errorf("no record found")
		}

		return json.Unmarshal(value, v)
	})
}

// Delete deletes an element from the bolt database
func Delete(bucket, id string) error {

	d := db()
	defer d.Close()

	//
	return d.Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return bucket.Delete([]byte(id))
	})
}

// Keys ...
func Keys(bucket string) (strArr []string, err error) {

	d := db()
	defer d.Close()

	//
	err = d.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return nil
		}

		//
		return bucket.ForEach(func(k, v []byte) error {
			strArr = append(strArr, string(k))
			return nil
		})
	})

	return
}
