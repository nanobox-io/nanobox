package update

import (
	"fmt"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

// 
const checkFrequency = (168 * time.Hour)

// Check checks to see if there is an update available for the nanobox CLI
func Check() {

	// load the update model
	update, err := models.LoadUpdate()
	if err != nil || update.CurrentVersion == "" {
		populateUpdate(update)
	}

	// return early if it's not time to check yet
	if !checkable(update) {
		return
	}

	// if the md5s don't match and it's 'time' for an update (14 days) inform the
	// user that updates are available
	if update.CurrentVersion != RemoveMd5() {

		//
		fmt.Printf(`
------------------------------------------------
Hey! A newer version of nanobox is available.
Run the following command to update:

$ nanobox-update
------------------------------------------------
`)

		// renew the update
	}

	update.LastCheckAt = time.Now()
	if err := update.Save(); err != nil {
		lumber.Error("update:update.Save(): %s", err)
	}
}

// Expired determines if the update has expired based on the expirationDate
// provided
func checkable(update *models.Update) bool {
	return time.Since(update.LastCheckAt) >= checkFrequency
}