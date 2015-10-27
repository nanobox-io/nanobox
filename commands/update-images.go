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
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"github.com/spf13/cobra"
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
	go Mist.Stream([]string{"log", "deploy"}, Mist.PrintLogStream)

	// listen for status updates
	errch := make(chan error)
	go func() {
		errch <- Mist.Listen([]string{"job", "imageupdate"}, Mist.ImageUpdates)
	}()

	// run an image update
	if err := Server.Update(""); err != nil {
		config.Fatal("[commands/update-images] server.Update() failed - ", err.Error())
	}

	// wait for a status update (blocking)
	err := <-errch

	//
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// PostRun: halt
}
