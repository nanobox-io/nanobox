// Package counter ...
package counter

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Increment will increment a counter by ID
func Increment(id string) (int, error) {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	counter, _ := models.FindCounterByID(id)

	// set the ID in case the counter was empty
	counter.ID = id

	//
	if counter.Count < 0 {
		counter.Count = 0
	}

	counter.Count = counter.Count + 1

	//
	err := counter.Save()

	return counter.Count, err
}

// Decrement will decrement a counter by ID
func Decrement(id string) (int, error) {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	counter, _ := models.FindCounterByID(id)

	// set the ID in case the counter was empty
	counter.ID = id

	//
	if counter.Count <= 0 {
		counter.Count = 0
	} else {
		counter.Count = counter.Count - 1
	}

	//
	err := counter.Save()

	return counter.Count, err
}

// Get returns the current counter value by ID
func Get(id string) (int, error) {

	counter, err := models.FindCounterByID(id)

	return counter.Count, err
}

// Reset resets a counter by ID; we can simply delete the counter model
func Reset(id string) error {
	locker.GlobalLock()
	defer locker.GlobalUnlock()

	counter := models.Counter{ID: id}
	return counter.Delete()
}
