package models

import (
	"fmt"
)

type Console struct {
	ID          string
	ContainerID string
}

// Save persists the Console to the database
func (c *Console) Save() error {

	if err := put("consoles", c.ID, c); err != nil {
		return fmt.Errorf("failed to save dev: %s", err.Error())
	}

	return nil
}

// Delete deletes the dev record from the database
func (c *Console) Delete() error {

	if err := destroy("consoles", c.ID); err != nil {
		return fmt.Errorf("failed to delete dev: %s", err.Error())
	}

	// clear the current entry
	c = nil

	return nil
}

// AllConsoles loads all of the Consoles in the database
func AllConsoles() ([]*Console, error) {
	// list of console to return
	console := []*Console{}

	// fetch all of the keys
	keys, err := keys("consoles")
	if err != nil {
		return nil, fmt.Errorf("failed to load dev keys: %s", err.Error())
	}

	// iterate over the keys and load each dev
	for _, key := range keys {
		dev, err := findConsoleByID(key, false)
		if err != nil {
			return nil, fmt.Errorf("failed to load dev record (%s): %s", key, err.Error())
		}

		console = append(console, dev)
	}

	return console, nil
}

// DeleteAllConsoles deletes all console
func DeleteAllConsoles() error {
	if err := truncate("consoles"); err != nil {
		return fmt.Errorf("failed to delete all console: %s", err.Error())
	}

	return nil
}

// findConsoleByID finds a dev by the ID and
// optionally errors if the record doesn't already exist
func findConsoleByID(ID string, mustExist bool) (*Console, error) {

	// create the exec object just incase we cant find one
	dev := &Console{}

	if err := get("consoles", ID, &dev); err != nil {

		// don't return an error if the record doesn't exist
		if !mustExist && err.Error() == "no record found" {
			return dev, nil
		}

		return dev, fmt.Errorf("failed to load dev: %s", err.Error())
	}

	return dev, nil
}
