package models

import (
	"fmt"
)

// Provider ...
type Provider struct {
	HostIP  string // the ip the host knows about
	MountIP string // the ip we reserved for mounting
}

// Save persists the Provider to the database
func (a *Provider) Save() error {

	// Since there is only ever a single provider value, we'll use the registry
	if err := put("registry", "provider", a); err != nil {
		return fmt.Errorf("failed to save provider: %s", err.Error())
	}

	return nil
}

// Delete deletes the provider record from the database
func (a *Provider) Delete() error {

	// Since there is only ever a single provider value, we'll use the registry
	if err := destroy("registry", "provider"); err != nil {
		return fmt.Errorf("failed to delete provider: %s", err.Error())
	}

	// clear the current entry
	a = nil

	return nil
}

// LoadProvider loads the provider entry
func LoadProvider() (*Provider, error) {
	provider := &Provider{}

	if err := get("registry", "provider", &provider); err != nil {
		return provider, fmt.Errorf("failed to load provider: %s", err.Error())
	}

	return provider, nil
}
