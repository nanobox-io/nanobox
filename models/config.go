package models

import (
	"fmt"
	"net"
)

// Config ...
type Config struct {
	Provider string `json:"provider"`

	// required for docker-machine
	MountType      string `json:"mount-type"`
	NetfsMountOpts string `json:"netfs-mount-opts"`
	CPUs           int    `json:"cpus"`
	RAM            int    `json:"ram"`
	Disk           int    `json:"disk"`

	// ip address spaces
	ExternalNetworkSpace      string `json:"external-network-space"`
	DockerMachineNetworkSpace string `json:"docker-machine-network-space"`
	NativeNetworkSpace        string `json:"native-network-space"`

	LockPort int `json:"lock-port"`
}

// Save persists the Config to the database
func (c *Config) Save() error {
	// make sure the information in is valid
	c.makeValid()

	// Since there is only ever a single Config value, we'll use the registry
	if err := put("registry", "Config", c); err != nil {
		return fmt.Errorf("failed to save Config: %s", err.Error())
	}

	return nil
}

// set reasonable default values for all necessary values
func (c *Config) makeValid() {
	if c.Provider != "native" && c.Provider != "docker-machine" {
		c.Provider = "docker-machine"
	}

	if c.MountType != "native" && c.MountType != "netfs" {
		c.MountType = "native"
	}

	if c.CPUs == 0 {
		c.CPUs = 1
	}

	if c.RAM == 0 {
		c.RAM = 1
	}

	if c.Disk < 20480 {
		c.Disk = 20480
	}

	if _, _, err := net.ParseCIDR(c.ExternalNetworkSpace); c.ExternalNetworkSpace == "" || err != nil {
		c.ExternalNetworkSpace = "192.168.99.50/24"
	}

	if _, _, err := net.ParseCIDR(c.DockerMachineNetworkSpace); c.DockerMachineNetworkSpace == "" || err != nil {
		c.DockerMachineNetworkSpace = "172.19.0.1/16"
	}

	if _, _, err := net.ParseCIDR(c.NativeNetworkSpace); c.NativeNetworkSpace == "" || err != nil {
		c.NativeNetworkSpace = "172.18.0.1/16"
	}

	if c.LockPort == 0 {
		c.LockPort = 12345
	}

}

// Delete deletes the Config record from the database
func (c *Config) Delete() error {

	// Since there is only ever a single Config value, we'll use the registry
	if err := destroy("registry", "Config"); err != nil {
		return fmt.Errorf("failed to delete Config: %s", err.Error())
	}

	// clear the current entry
	c = nil

	return nil
}

// LoadConfig loads the Config entry
func LoadConfig() (*Config, error) {
	c := &Config{}
	c.makeValid()
	if err := get("registry", "Config", &c); err != nil {
		return c, fmt.Errorf("failed to load Config: %s", err.Error())
	}

	return c, nil
}
