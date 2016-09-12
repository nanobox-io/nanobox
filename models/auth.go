package models

import (
	"fmt"
)

// Auth ...
type Auth struct {
	Endpoint string // nanobox, bonesalt, dev, sim
	Key      string // api_token from dashboard
}

// determines if the auth record is new
func (a *Auth) IsNew() bool {
	return a.Key == ""
}

// Save persists the Auth to the database
func (a *Auth) Save() error {

	// Since there is only ever a single auth value, we'll use the registry
	if err := put("auths", a.Endpoint, a); err != nil {
		return fmt.Errorf("failed to save auth: %s", err.Error())
	}

	return nil
}

// Delete deletes the auth record from the database
func (a *Auth) Delete() error {

	return DeleteAuth(a.Endpoint)
}

// LoadAuth loads the default (nanobox) auth entry
func LoadAuth() (*Auth, error) {
	auth := &Auth{
		Endpoint: "nanobox",
	}

	if err := get("auths", auth.Endpoint, &auth); err != nil {
		return auth, fmt.Errorf("failed to load auth: %s", err.Error())
	}

	return auth, nil
}

// loads an auth by a specific endpoint
func LoadAuthByEndpoint(endpoint string) (*Auth, error) {
	auth := &Auth{
		Endpoint: endpoint,
	}

	if err := get("auths", endpoint, &auth); err != nil {
		return auth, fmt.Errorf("failed to load auth: %s", err.Error())
	}

	return auth, nil
}

// DeleteAuth ...
func DeleteAuth(endpoint string) error {

	if err := destroy("auths", endpoint); err != nil {
		return fmt.Errorf("failed to delete auth: %s", err.Error())
	}

	return nil
}
