package counter

import (
  "github.com/nanobox-io/nanobox/models"
  "github.com/nanobox-io/nanobox/util/data"
  "github.com/nanobox-io/nanobox/util/locker"
)

// Increment will increment a counter by ID
func Increment(ID string) (int, error) {

  // aquire a lock to ensure we're atomically updating
  if err := locker.GlobalLock(); err != nil {
    return 0, err
  }

  counter := models.Counter{}

  data.Get("counters", ID, &counter)

  if counter.Count < 0 {
    counter.Count = 0
  }

  counter.Count = counter.Count + 1

  if err := data.Put("counters", ID, counter); err != nil {
    return 0, err
  }

  // release the lock
  if err := locker.GlobalUnlock(); err != nil {
    return 0, err
  }

  return counter.Count, nil
}

// DecrementCounter will decrement a counter by ID
func Decrement(ID string) (int, error) {

  // aquire a lock to ensure we're atomically updating
  if err := locker.GlobalLock(); err != nil {
    return 0, err
  }

  counter := models.Counter{}

  data.Get("counters", ID, &counter)

  if counter.Count <= 0 {
    counter.Count = 0
  } else {
    counter.Count = counter.Count - 1
  }

  if err := data.Put("counters", ID, counter); err != nil {
    return 0, err
  }

  // release the lock
  if err := locker.GlobalUnlock(); err != nil {
    return 0, err
  }

  return counter.Count, nil
}

// Get returns the current counter value by ID
func Get(ID string) (int, error) {
  counter := models.Counter{}

  if err := data.Get("counters", ID, &counter); err != nil {
    return 0, err
  }

  return counter.Count, nil
}

// Reset resets a counter by ID
func Reset(ID string) error {
  // we can simply delete the counter model

  // aquire a lock to ensure we're atomically updating
  if err := locker.GlobalLock(); err != nil {
    return err
  }

  // we don't care if it errors because that means it didn't exist
  data.Delete("counters", ID)

  // release the lock
  if err := locker.GlobalUnlock(); err != nil {
    return err
  }

  return nil
}
