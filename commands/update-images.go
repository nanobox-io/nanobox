// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/server"
	"github.com/nanobox-io/nanobox-cli/util/server/mist"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var updateImagesCmd = &cobra.Command{
	Use:   "update-images",
	Short: "Updates the nanobox docker images",
	Long:  ``,

	PreRun:  boot,
	Run:     updateImages,
	PostRun: halt,
}

// updateImages
func updateImages(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	fmt.Printf(stylish.Bullet("Updating nanobox docker images..."))

	// stream update output
	go mist.Stream([]string{"log", "deploy"}, mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- mist.Listen([]string{"job", "imageupdate"}, mist.ImageUpdates)
	}()

	// run an image update
	if err := server.Update(""); err != nil {
		config.Fatal("[commands/imagesUpdate] failed - ", err.Error())
	}

	// wait for a status update (blocking)
	err := <-errch

	switch {

	//
	case err == nil:

	//
	case err != nil:
		fmt.Printf(err.Error())
	}

	// PostRun: halt
}
