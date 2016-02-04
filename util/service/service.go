//
package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nanobox-io/nanobox/config"
	fileutil "github.com/nanobox-io/nanobox/util/file"
)

// RemountLocal simply calls MountLocal() but only returns the error as the
// mounts name and dir are not important in this instance.
func RemountLocal() (err error) {
	_, _, err = MountLocal()
	return
}

// MountLocal creates a local mount (~/.nanobox/apps/<app>/<service>/<mount>)
func MountLocal() (mountName, mountDir string, err error) {

	// parse the boxfile and see if there is an service declared; if none is declared
	// simply return.
	servicePath := config.ParseBoxfile().Build.Service
	if servicePath == "" {
		return
	}

	//
	mountName = filepath.Base(servicePath)

	// if no local service is found return since there is nothing more to do here;
	// when an service is specified but not found, it's assumed that the desired
	// service exists on nanobox.io
	if _, err = os.Stat(servicePath); err != nil {
		return
	}

	//
	servicefile := filepath.Join(servicePath, "Servicefile")

	// ensure there is an servicefile at the service location
	if _, err = os.Stat(servicefile); err != nil {
		err = fmt.Errorf("No servicefile found at '%v', Exiting...\n", servicePath)
		return
	}

	// if there is an servicefile attempt to parse it to get any additional build
	// files for mounting
	files := []string{"./bin", "./Servicefile", "./meta.json"}
	if err = config.ParseConfig(servicefile, files); err != nil {
		err = fmt.Errorf("Nanobox failed to parse your Servicefile. Please ensure it is valid YAML and try again.\n")
		return
	}

	// directory to mount
	mountDir = filepath.Join(config.AppDir, mountName)

	// if the servicefile parses successfully create the mount only if it doesn't
	// already exist
	if err = os.MkdirAll(mountDir, 0755); err != nil {
		return
	}

	var abs string

	//
	abs, err = filepath.Abs(servicePath)
	if err != nil {
		return
	}

	// pull the remaining service files into the mount
	for _, files := range files {

		path := filepath.Join(abs, files)

		// just skip any files that aren't found; any required files will be
		// caught before publishing, here it doesn't matter
		if _, err := os.Stat(path); err != nil {
			continue
		}

		// copy service file into mount
		if err = fileutil.Copy(path, mountDir); err != nil {
			return
		}
	}

	return
}
