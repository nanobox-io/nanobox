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

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,

	Run: Update,
}

// Update
func Update(ccmd *cobra.Command, args []string) {

	// install the nanobox vagrant image only if it isn't already available
	if !vagrant.HaveImage() {
		fmt.Printf(stylish.Bullet("Installing nanobox image..."))

		// install the nanobox vagrant image
		if err := vagrant.Install(); err != nil {
			config.Fatal("[commands/boxInstall] failed - ", err.Error())
		}
	}

	//
	match, err := util.MD5sMatch(config.Root+"/nanobox-boot2docker.md5", "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5")
	if err != nil {
		config.Fatal("[commands/boxUpdate] failed - ", err.Error())
	}

	// if the local md5 does not match the remote md5 download the newest nanobox
	// image
	if !match {
		fmt.Printf(stylish.Bullet("Updating nanobox image..."))

		// update the nanobox vagrant image
		if err := vagrant.Update(); err != nil {
			config.Fatal("[commands/boxUpdate] failed - ", err.Error())
		}
	}
}
