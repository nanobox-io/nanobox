package models

import (
	"fmt"
)

// Counter provides a generic dataset to increment and decrement counters across
// the entire counterlication with varying levels of granularity
type Counter struct {
	ID    string //
	Count int    //
}

// IsNew returns true if the Counter hasn't been created yet
func (c *Counter) IsNew() bool {
	return c.ID == ""
}

// Save persists the Counter to the database
func (c *Counter) Save() error {
	
	if err := put("counters", c.ID, c); err != nil {
		return fmt.Errorf("failed to save counter: %s", err.Error())
	}
	
	return nil
}

// Delete deletes the counter record from the database
func (c *Counter) Delete() error {
	
	if err := delete("counters", c.ID); err != nil {
		return fmt.Errorf("failed to delete counter: %s", err.Error())
	}
	
	// clear the current entry
	c = nil
	
	return nil
}

// FindCounterByID finds an counter by an ID
func FindCounterByID(ID string) (Counter, error) {
	return findCounterByID(ID, false)
}

// MustFindCounterByID finds a counter by an ID and 
// returns an error if the counter doesn't exist
func MustFindCounterByID(ID string) (Counter, error) {
	return findCounterByID(ID, true)
}

// AllCounters loads all of the Counters in the database
func AllCounters() ([]Counter, error) {
	// list of counters to return
	counters := []Counter{}
	
	// fetch all of the keys
	keys, err := keys("counters")
	if err != nil {
		return nil, fmt.Errorf("failed to load counter keys: %s", err.Error())
	}
	
	// iterate over the keys and load each counter
	for _, key := range keys {
		counter, err := FindCounterByID(key)
		if err != nil {
			return nil, fmt.Errorf("failed to load counter record (%s): %s", key, err.Error())
		}
		
		counters = append(counters, counter)
	}
	
	return counters, nil
}

// DeleteAllCounters deletes all counters
func DeleteAllCounters() error {
	if err := truncate("counters"); err != nil {
		return fmt.Errorf("failed to delete all counters: %s", err.Error())
	}
	
	return nil
}

// findCounterByID finds a counter by the ID and 
// optionally errors if the record doesn't already exist
func findCounterByID(ID string, mustExist bool) (Counter, error) {
	
	counter := Counter{}
	
	if err := get("counters", ID, &counter); err != nil {
		
		// don't return an error if the record doesn't exist
		if !mustExist && err.Error() == "no record found" {
			return counter, nil
		}
		
		return counter, fmt.Errorf("failed to load counter: %s", err.Error())
	}
	
	return counter, nil
}
