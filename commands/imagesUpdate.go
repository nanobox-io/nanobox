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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var imagesUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the nanobox docker images",
	Long:  ``,

	PreRun:  bootVM,
	Run:     imagesUpdate,
	PostRun: saveVM,
}

// imagesUpdate
func imagesUpdate(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	fmt.Printf(stylish.Bullet("Updating nanobox docker images..."))

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
	case "complete":
		break

	case "errored":
		fmt.Printf("Nanobox failed to update docker images")
	}

	// PostRun: saveVM
}
