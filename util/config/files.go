package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util"
)

// Boxfile ...
func Boxfile() string {
	return filepath.ToSlash(filepath.Join(LocalDir(), "boxfile.yml"))
}

// configFile returns the path/to/config.yml file or creates one if it doesn't
// already exist; when created the file is populated with resonable defaults.
// This is mainly used as an example for users to see what a config file can look
// like. once they create their own we'll use that with any remaining defaults
// pulled from viper (see ./config.go)
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
		lumber.Fatal("[util/config/files] os.Create() failed", err.Error())
	}
	defer f.Close()

	//
	contents := fmt.Sprintf(`
# provider configuration options
provider: "docker_machine" # the name of the provider to use

# mount type (native|netfs)
mount-type: native

token: %s
  `, util.RandomString(30))

	// populate the config.yml with reasonable defaults
	io.WriteString(f, contents)

	return
}
