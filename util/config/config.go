package config

import (
	"os"
	"path/filepath"

	"github.com/nanobox-io/nanobox/util"
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
			return nil
		}
	}

	return config
}

// ParseConfig will load the config file and parse it with viper
func parseConfig() error {

	// initilize a viper parser
	config = viper.New()

	// set sensible defaults
	// network space
	config.SetDefault("external-network-space", "192.168.99.50/24")
	config.SetDefault("internal-network-space", "192.168.0.50/16")

	// default provider
	config.SetDefault("provider", "docker_machine")

	// vm properties for vm providers
	config.SetDefault("vm.cpus", 2)
	config.SetDefault("vm.cpu-cap", 50)
	config.SetDefault("vm.ram", 1024)
	config.SetDefault("vm.mount", "native")

	// odin access
	config.SetDefault("production_host", "api.nanobox.io/v1/")

	// no sense parsing the config file if it doesn't exist
	if configExists() {

		// parse config file
		configFile := filepath.Join(util.GlobalDir(), "config.yml")
		config.SetConfigFile(configFile)

		// merge with defaults
		if err := config.MergeInConfig(); err != nil {
			return err
		}
	}

	return nil
}

// configExists checks to see if the config file actually exists
func configExists() bool {

	// we simply stat the file and if there are no errors the file exists
	configFile := filepath.Join(util.GlobalDir(), "config.yml")
	if _, err := os.Stat(configFile); err == nil {
		return true
	}

	return false
}
