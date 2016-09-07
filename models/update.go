package models

import (
	"fmt"
	"time"
)

// Update ...
type Update struct {
	LastUpdatedAt time.Time
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

// Renew sets the LastUpdatedAt of the Update to time.Now() effectively "renewing"
// the expiration
func (u *Update) Renew() error {
	u.LastUpdatedAt = time.Now()
	return u.Save()
}

// Expired determines if the update has expired based on the expirationDate
// provided
func (u *Update) Expired(expirationDate float64) bool {
	return time.Since(u.LastUpdatedAt).Hours() >= expirationDate
}
