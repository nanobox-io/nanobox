// Package update ...
package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kardianos/osext"
	"github.com/nanobox-io/nanobox/models"
	cryptoutil "github.com/sdomino/go-util/crypto"
)

// EXPIREAFTER is the time in hours after which we want to check for updates (168 hours or 7 days)
const EXPIREAFTER = 168

// Check checks to see if there is an update available for the nanobox CLI
func Check() error {

	// load the update model
	update, _ := models.LoadUpdate()

	// return early if it's not time to update yet
	if !update.Expired(EXPIREAFTER) {
		return nil
	}

	//
	updatable, err := updatable()
	if err != nil {
		return fmt.Errorf("Nanobox was unable to determine if updates are available - %s", err.Error())
	}

	// if the md5s don't match and it's 'time' for an update (14 days) inform the
	// user that updates are available
	if updatable {

		//
		fmt.Printf(`
------------------------------------------------
Hey! A newer version of nanobox is available.
Run the following command to update:

$ nanobox-update
------------------------------------------------
`)

		// renew the update
		return update.Renew()
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
