package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// "github.com/jcelliott/lumber"
	// "github.com/nanobox-io/nanobox/util"
)

type SetupConf struct {
	Provider string
	Mount    string
	CPUs     int
	RAM      int
}

// Boxfile ...
func Boxfile() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}


func ConfigExists() bool {
	file := filepath.ToSlash(filepath.Join(GlobalDir(), "config.yml"))

	// if i can stat the file exists
	_, err := os.Stat(file)
	return err == nil
	
}

// configFile returns the path/to/config.yml file or creates one if it doesn't
// already exist; when created the file is populated with resonable defaults.
// This is mainly used as an example for users to see what a config file can look
// like. once they create their own we'll use that with any remaining defaults
// pulled from viper (see ./config.go)
func ConfigFile(setup *SetupConf) (file string) {
	//
	file = filepath.ToSlash(filepath.Join(GlobalDir(), "config.yml"))

	// return the filepath if it's already created...
	if _, err := os.Stat(file); err == nil && setup == nil {
		return
	}

	// 
	if setup == nil {
		setup = &SetupConf{
			Provider: "docker-machine",
			Mount: "native",
			CPUs: 1,
			RAM: 1,
		}
	}
	// ...otherwise create the file

	//
	contents := fmt.Sprintf(`

# provider configuration options
provider: "%s" # the name of the provider to use

# This next section is used by the docker-machine setup only

# mount type (native|netfs)
mount-type: %s

# number of cpus you want docker-machine to have access to
cpus: %d

# number of gigabytes of ram you want docker-machine to use
ram: %d
`, setup.Provider, setup.Mount, setup.CPUs, setup.RAM)

	// populate the config.yml with reasonable defaults
	ioutil.WriteFile(file, []byte(contents), 0666)

	return
}
