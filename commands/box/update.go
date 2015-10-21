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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,

	Run: Update,
}

// Update
func Update(ccmd *cobra.Command, args []string) {

	if err := checkInstall(); err != nil {
		Config.Fatal("[commands/boxInstall] failed - ", err.Error())
	}

	//
	match, err := Util.MD5sMatch(Config.Root()+"/nanobox-boot2docker.md5", "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5")
	if err != nil {
		Config.Fatal("[commands/boxUpdate] failed - ", err.Error())
	}

	// if the local md5 does not match the remote md5 download the newest nanobox
	// image
	if !match {
		fmt.Printf(stylish.Bullet("Updating nanobox image..."))

		// update the nanobox vagrant image
		if err := Vagrant.Update(); err != nil {
			Config.Fatal("[commands/boxUpdate] failed - ", err.Error())
		}
	}
}
