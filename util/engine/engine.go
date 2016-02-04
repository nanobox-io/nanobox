//
package engine

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

// MountLocal creates a local mount (~/.nanobox/apps/<app>/<engine>/<mount>)
func MountLocal() (mountName, mountDir string, err error) {

	// parse the boxfile and see if there is an engine declared; if none is declared
	// simply return.
	enginePath := config.ParseBoxfile().Build.Engine
	if enginePath == "" {
		return
	}

	//
	mountName = filepath.Base(enginePath)

	// if no local engine is found return since there is nothing more to do here;
	// when an engine is specified but not found, it's assumed that the desired
	// engine exists on nanobox.io
	if _, err = os.Stat(enginePath); err != nil {
		return
	}

	//
	enginefile := filepath.Join(enginePath, "Enginefile")

	// ensure there is an enginefile at the engine location
	if _, err = os.Stat(enginefile); err != nil {
		err = fmt.Errorf("No enginefile found at '%v', Exiting...\n", enginePath)
		return
	}

	// if there is an enginefile attempt to parse it to get any additional build
	// files for mounting
	files := []string{"./bin", "./Enginefile", "./meta.json"}
	if err = config.ParseConfig(enginefile, files); err != nil {
		err = fmt.Errorf("Nanobox failed to parse your Enginefile. Please ensure it is valid YAML and try again.\n")
		return
	}

	// directory to mount
	mountDir = filepath.Join(config.AppDir, mountName)

	// if the enginefile parses successfully create the mount only if it doesn't
	// already exist
	if err = os.MkdirAll(mountDir, 0755); err != nil {
		return
	}

	var abs string

	//
	abs, err = filepath.Abs(enginePath)
	if err != nil {
		return
	}

	// pull the remaining engine files into the mount
	for _, files := range files {

		path := filepath.Join(abs, files)

		// just skip any files that aren't found; any required files will be
		// caught before publishing, here it doesn't matter
		if _, err := os.Stat(path); err != nil {
			continue
		}

		// copy engine file into mount
		if err = fileutil.Copy(path, mountDir); err != nil {
			return
		}
	}

	return
}
