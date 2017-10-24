package update

import (
	"fmt"
	"strings"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
)

// Check for updates to nanobox every other day
const checkFrequency = (48 * time.Hour)

// Check checks to see if there is an update available for the nanobox CLI
func Check() {
	// load the update model
	updateInfo, err := models.LoadUpdate()
	if err != nil {
		if strings.Contains(err.Error(), "no record found") {
			checkTomorrow(&models.Update{})
			return
		}
		lumber.Error("update:models.LoadUpdate(): %s", err)
		return
	}

	// return early if it's not time to check yet
	if !checkable(updateInfo) {
		return
	}

	latest := latestVersion()
	if latest == "" {
		checkTomorrow(updateInfo)
		return
	}
	lumber.Debug("CurrVers: %s\nLatest:   %s\n", models.VersionString(), latest)

	// if the versions don't match and it's 'time' for an update (14 days) inform the
	// user that updates are available
	if models.VersionString() != latest {
		//
		fmt.Printf(`
------------------------------------------------
Hey! A newer version of nanobox is available.

  %s

Run the following command to update:

$ nanobox-update
------------------------------------------------
`, latest)

	}

	// renew the update last checked time
	updateInfo.LastCheckAt = time.Now()
	if err := updateInfo.Save(); err != nil {
		lumber.Error("update:updateInfo.Save(): %s", err)
	}
}

// Checkable determines if the update has expired based on the check frequency
func checkable(updateInfo *models.Update) bool {
	return time.Since(updateInfo.LastCheckAt) >= checkFrequency
}

// CheckTomorrow updates the last checked at time to a day later so nanobox will
// check for updates tommorrow. This is only called if the check failed, likely
// due to a network error.
func checkTomorrow(updateInfo *models.Update) error {
	updateInfo.LastCheckAt = updateInfo.LastCheckAt.Add(24 * time.Hour)
	return updateInfo.Save()
}
