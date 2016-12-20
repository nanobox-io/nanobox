// Package config ...
package config

import (
	"github.com/jcelliott/lumber"
	"github.com/spf13/viper"
)

//
var config *viper.Viper

// Get fetches a generic value from viper config
func Get(key string) interface{} {

	// parse the config file if it's not already
	if config == nil {
		if err := parseConfig(); err != nil {
			return nil
		}
	}

	return config.Get(key)
}

// Viper returns the viper config object
func Viper() *viper.Viper {

	// parse the config file if it's not already
	if config == nil {
		if err := parseConfig(); err != nil {
			lumber.Error("config:Viper():parseConfig(): %s", err.Error())
		}
	}

	return config
}

// ParseConfig will load the config file and parse it with viper
func parseConfig() error {

	// initilize a viper parser
	config = viper.New()

	// set sensible defaults

	// network spaces
	config.SetDefault("native-network-space", "172.18.0.10/16")

	config.SetDefault("external-network-space", "192.168.99.50/24")
	config.SetDefault("docker-machine-network-space", "172.19.0.10/16")

	// configurable options (via ~.nanobox/config.yml); these defaults are set here
	// incase for some reason an incomplete or invalid config.yml is parsed, nanobox
	// will have values to fall back on
	config.SetDefault("provider", "docker-machine")
	config.SetDefault("mount-type", "native")

	// parse config file; we attempt to parse the config.yml and pull out any values
	// that the user has provided (or one is generated with defaults; see ./files.go)
	config.SetConfigFile(ConfigFile(nil))

	// merge with defaults
	if err := config.MergeInConfig(); err != nil {
		return err
	}

	return nil
}
