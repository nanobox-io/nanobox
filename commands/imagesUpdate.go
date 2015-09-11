// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var imagesUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the nanobox docker images",
	Long: `
Description:
  Updates the nanobox docker images`,

	Run: imagesUpdate,
}

// imagesUpdate
func imagesUpdate(ccmd *cobra.Command, args []string) {
	fmt.Printf(stylish.Bullet("Updating nanobox docker images..."))

	//
	update := util.Sync{
		Model:   "imageupdate",
		Path:    fmt.Sprintf("http://%s/image-update", config.ServerURI),
		Verbose: fVerbose,
	}

	//
	update.Run(args)

	//
	switch update.Status {

	// complete
	case "complete":
		fmt.Printf(stylish.Bullet("Update complete"))

	// if the bootstrap fails the server should handle the message. If not, this can
	// be re-enabled
	case "errored":
		// fmt.Printf(stylish.Error("Bootstrap failed", "Your app failed to bootstrap"))
	}
}
