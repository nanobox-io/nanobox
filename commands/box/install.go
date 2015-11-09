// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

//
package box

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: Install,
}

// Install
func Install(ccmd *cobra.Command, args []string) {

	// if the nanobox-boot2docker.box is not installed, download and install it
	if err := checkInstall(); err != nil {
		Config.Fatal("[commands/box/install] checkInstall() failed - ", err.Error())
	}
}

// checkInstall determines if the nanobox-boot2docker.box needs to be downloaded
func checkInstall() (err error) {

	// download nanobox-boot2docker.box only if it isn't already available
	if !Vagrant.HaveImage() {
		fmt.Printf(stylish.Bullet("Installing nanobox image..."))

		// 'install' nanobox-boot2docker.box
		err = Vagrant.Install()
	}
	return
}
