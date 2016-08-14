package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/nanobox-io/nanobox/util/config"
)

var (
  // path to database
  DB = filepath.ToSlash(filepath.Join(config.GlobalDir(), "data.db"))
)

// db opens a boltDB connection
func db() (*bolt.DB, error) {

	boltDB, err := bolt.Open(DB, 0666, nil)
	if err != nil {
    return nil, fmt.Errorf("unable to open database file (%s): %s", DB, err.Error())
  }

	return boltDB, nil
}

// put inserts or updates an element into the bolt database
func put(bucket, id string, v interface{}) error {

  // open the database
  db, err := db()
  if err != nil {
    return fmt.Errorf("unable to initialize database driver: %s ", err.Error())
  }
  
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {

		// Create a bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
      return fmt.Errorf("unable to create a database bucket: %s", err.Error())
		}

		// Marshal the value into a JSON blob
		bytes, err := json.Marshal(v)
		if err != nil {
      return fmt.Errorf("failed to encode database record: %s", err.Error())
		}
    
    // Write the entry
		if err := bucket.Put([]byte(id), bytes); err != nil {
      return fmt.Errorf("failed to write entry: %s", err.Error())
    }

		return nil
	})
}

// get retrieves an element from the database unlike the default behavior from
// boltdb this will return an error if you try getting something that doesnt exist
func get(bucket, id string, v interface{}) error {

	// open the database
  db, err := db()
  if err != nil {
    return fmt.Errorf("unable to initialize database driver: %s ", err.Error())
  }
  
	defer db.Close()

  // Read value back in a read-only transaction.
	return db.View(func(tx *bolt.Tx) error {

		// Establish the table (bucket)
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return fmt.Errorf("no record found")
		}

		// Fetch the value
		value := bucket.Get([]byte(id))
		if value == nil || len(value) == 0 {
			return fmt.Errorf("no record found")
		}

    if err := json.Unmarshal(value, v); err != nil {
      return fmt.Errorf("failed to decode database record: %s", err.Error())
    }

		return nil
	})
}

// delete deletes an element from the bolt database
func delete(bucket, id string) error {

	// open the database
  db, err := db()
  if err != nil {
    return fmt.Errorf("unable to initialize database driver: %s ", err.Error())
  }
  
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
    
		// Create a bucket.
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
      return fmt.Errorf("unable to create database bucket: %s", err.Error())
		}

    // Delete the record from the table (bucket)
    if err := bucket.Delete([]byte(id)); err != nil {
      return fmt.Errorf("failed to delete database record: %s", err.Error())
    }
    
		return nil
	})
}

// keys returns a list of keys in a table (bucket)
func keys(bucket string) (keys []string, err error) {

	// open the database
  db, err := db()
  if err != nil {
    return nil, fmt.Errorf("unable to initialize database driver: %s ", err.Error())
  }
  
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
    
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return nil
		}

		return bucket.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})

	return
}

// get all the elements in a bucket and place them into an interface
// this will not work if v is not an array
func getAll(bucket string, v interface{}) error {
	// get the keys
	keys, err := keys(bucket)
	if err != nil {
		return fmt.Errorf("unable to get keys: %s", err)
	}

	// start making a array so we can marshel the all the 
	// elements into it
	elements := [][]byte{}

	// open the database 
  db, err := db()
  if err != nil {
    return fmt.Errorf("unable to initialize database driver: %s ", err.Error())
  }
  
	defer db.Close()

  // Read value back in a read-only transaction.
	err = db.View(func(tx *bolt.Tx) error {

		// Establish the table (bucket)
		bucket := tx.Bucket([]byte(bucket))
		if bucket == nil {
			return fmt.Errorf("no record found")
		}

		// Fetch the values and append them to the elements array
		for _, key := range keys {
			value := bucket.Get([]byte(key))
			if value == nil || len(value) == 0 {
				return fmt.Errorf("no record found")
			}
			elements = append(elements, value)
		}

		return nil
	})

	if err != nil  {
		if err.Error() == "no record found" {
			return nil
		}
		return fmt.Errorf("unable to load a key from the keys list (%s): %s", bucket, err)
	}

	// combine the json elements into a json array
	combinedElements := []byte(fmt.Sprintf("[%s]", bytes.Join(elements, []byte{','})))

	// unmarshel the data into the interface
  if err := json.Unmarshal(combinedElements, v); err != nil {
    return fmt.Errorf("failed to decode database record: %s", err.Error())
  }
  
  return nil
}

// truncate deletes a bucket and all entries
func truncate(bucket string) error {
	// fetch all the keys
	keys, err := keys(bucket)
	if err != nil {
		return fmt.Errorf("failed to fetch keys from bucket: %s", err.Error())
	}
	
	// delete the keys
	for _, key := range keys {
		if err := delete(bucket,key); err != nil {
			return fmt.Errorf("failed to delete entry by key (%s): %s", key, err.Error())
		}
	}
	
	return nil
}
