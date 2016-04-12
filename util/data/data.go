// Copyright (C) Pagoda Box, Inc - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
package data

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/boltdb/bolt"

	"github.com/nanobox-io/nanobox/util"
)

var boltDB *bolt.DB

// Establish a bolt.DB for future use.
func db() *bolt.DB {
	// reuse databse if one exists
	if boltDB != nil {
		return boltDB
	}
	var err error
	boltDB, err = bolt.Open(filepath.ToSlash(filepath.Join(util.GlobalDir(), "data.db")), 0666, nil)
	if err != nil {
		panic(err)
	}

	return boltDB
}

// Insert or update an element into the bolt database
func Put(bucket, id string, v interface{}) error {

	return db().Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		bucket.Put([]byte(id), bytes)
		return nil
	})
}

// retrieve an element from the database
// unlike the default behavior from boltdb this will return
// an error if you try getting something that doesnt exist
func Get(bucket, id string, v interface{}) error {
	// Read value back in a different read-only transaction.
	return db().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return fmt.Errorf("no record found")
		}
		value := bucket.Get([]byte(id))
		if value == nil || len(value) == 0 {
			return fmt.Errorf("no record found")
		}
		return json.Unmarshal(value, v)
	})
}

// delete an element from the bolt database
func Delete(bucket, id string) error {
	return db().Update(func(tx *bolt.Tx) error {
		// Create a bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}

		return bucket.Delete([]byte(id))

	})
}

func Keys(bucket string) ([]string, error) {
	strArr := []string{}
	err := db().View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))
		return bucket.ForEach(func(k, v []byte) error {
			strArr = append(strArr, string(k))
			return nil
		})
	})
	return strArr, err
}
