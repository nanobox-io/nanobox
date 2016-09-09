package models

import (
	"fmt"
	"time"
)

// Update ...
type Update struct {
	CurrentVersion string
	LastCheckAt    time.Time
	LastUpdatedAt  time.Time
}

// LoadUpdate loads the update entry
func LoadUpdate() (*Update, error) {
	update := &Update{}

	if err := get("registry", "update", &update); err != nil {
		return update, fmt.Errorf("failed to load update: %s", err.Error())
	}

	return update, nil
}

// Save persists the Update to the database
func (u *Update) Save() error {

	// Since there is only ever a single update value, we'll use the registry
	if err := put("registry", "update", u); err != nil {
		return fmt.Errorf("failed to save update: %s", err.Error())
	}

	return nil
}
