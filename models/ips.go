package models

import (
	"fmt"
	"net"
)

// IPs ...
type IPs []net.IP

// Save persists the IPs to the database
func (ips *IPs) Save() error {

	// Since there is only ever ips single auth value, we'll use the registry
	if err := put("registry", "ips", ips); err != nil {
		return fmt.Errorf("failed to save auth: %s", err.Error())
	}

	return nil
}

// Delete deletes the auth record from the database
func (ips *IPs) Delete() error {

	// Since there is only ever a single auth value, we'll use the registry
	if err := destroy("registry", "ips"); err != nil {
		return fmt.Errorf("failed to delete auth: %s", err.Error())
	}

	return nil
}

// LoadIPs loads the auth entry
func LoadIPs() (IPs, error) {
	ips := IPs{}

	if err := get("registry", "ips", &ips); err != nil {
		return ips, fmt.Errorf("failed to load ips: %s", err.Error())
	}

	return ips, nil
}
