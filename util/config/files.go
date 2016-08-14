package config

import (
	"io"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"
)

// Boxfile ...
func Boxfile() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}

// configFile creates a config.yml file (if one doesn't already exist) populated
// with resonable defaults; this is mainly used as an example for users to see
// what a config file can look like. once they create their own we'll use that
// with any remaining defaults pulled from viper (see ./config.go)
func configFile() (file string) {

	//
	file = filepath.ToSlash(filepath.Join(GlobalDir(), "config.yml"))

	// return the filepath if it's already created...
	if _, err := os.Stat(file); err == nil {
		return
	}

	// ...otherwise create the file
	f, err := os.Create(file)
	if err != nil {
		lumber.Fatal("[config/config] os.Create() failed", err.Error())
	}
	defer f.Close()

	//
	contents := `
# provider configuration options
provider: "docker_machine" # the name of the provider to use

# mount type (native|netfs)
mount-type: native

  `

	// populate the config.yml with reasonable defaults
	io.WriteString(f, contents)

	return
}
