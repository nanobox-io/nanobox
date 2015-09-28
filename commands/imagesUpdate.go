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

	PreRun: bootVM,
	Run:    imagesUpdate,
}

// imagesUpdate
func imagesUpdate(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	fmt.Printf(stylish.Bullet("Updating docker images..."))

	//
	update := util.Sync{
		Model:   "imageupdate",
		Path:    fmt.Sprintf("%s/image-update", config.ServerURL),
		Verbose: fVerbose,
	}

	//
	update.Run(args)

	//
	switch update.Status {

	//
	case "complete", "errored":
		break
	}
}
