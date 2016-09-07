// Package update ...
package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kardianos/osext"
	"github.com/nanobox-io/nanobox/util/config"
	cryptoutil "github.com/sdomino/go-util/crypto"
)

// Check checks to see if there is an update available for the nanobox CLI
func Check() error {

	//
	updatable, err := updatable()
	if err != nil {
		return fmt.Errorf("Nanobox was unable to determine if updates are available - %s", err.Error())
	}

	// stat the update file to get ModTime(); an error here means the file doesn't
	// exist, which is highly unlikely as this command creates it if it doesn't
	// exist, so we skip the error
	fi, _ := os.Stat(config.UpdateFile())

	// if the md5s don't match and it's 'time' for an update (14 days) inform the
	// user that updates are available
	if updatable && time.Since(fi.ModTime()).Hours() >= 336.0 {

		//
		fmt.Printf(`
# Update available
------------------------------------------------
A newer version of the nanobox CLI is available. We highly recommend updating at
your earliest convenience. Run the following command to update:

$ nanobox-update
------------------------------------------------`)

		// update the mod time on the updateFile file so we won't check for updates
		// again
		if err := os.Chtimes(config.UpdateFile(), time.Now(), time.Now()); err != nil {
			return err
		}
	}

	return nil
}

// updatable
func updatable() (bool, error) {

	//
	path, err := osext.Executable()
	if err != nil {
		return false, err
	}

	// check the md5 of the current executing cli against the remote md5;
	// os.Args[0] is used as the final interpolation to determine standard/dev versions
	match, err := cryptoutil.MD5Match(path, fmt.Sprintf("https://s3.amazonaws.com/tools.nanobox.io/cli/%s/%s/%s.md5", runtime.GOOS, runtime.GOARCH, filepath.Base(os.Args[0])))
	if err != nil {
		return false, err
	}

	// if the MD5's DONT match we want to update
	return !match, nil
}
