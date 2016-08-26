package models

import (
	"fmt"
)

type Auth struct {
	Key string //
}

// Save persists the Auth to the database
func (a *Auth) Save() error {

	// Since there is only ever a single auth value, we'll use the registry
	if err := put("registry", "auth", a); err != nil {
		return fmt.Errorf("failed to save auth: %s", err.Error())
	}

	return nil
}

// Delete deletes the auth record from the database
func (a *Auth) Delete() error {

	return DeleteAuth()
}

// LoadAuth loads the auth entry
func LoadAuth() (*Auth, error) {
	auth := &Auth{}

	if err := get("registry", "auth", &auth); err != nil {
		return auth, fmt.Errorf("failed to load auth: %s", err.Error())
	}

	return auth, nil
}

func DeleteAuth() error {
	// Since there is only ever a single auth value, we'll use the registry
	if err := destroy("registry", "auth"); err != nil {
		return fmt.Errorf("failed to delete auth: %s", err.Error())
	}

	return nil
}
