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

	var err error
	var match bool

	// if the nanobox-boot2docker.box is not installed, download and install it
	if err = checkInstall(); err != nil {
		Config.Fatal("[commands/box/update] checkInstall() failed - ", err.Error())
	}

	// ensure the local nanobox-boot2docker.box matches the remote one
	if match, err = Util.MD5sMatch(Config.Root()+"/nanobox-boot2docker.box", "https://s3.amazonaws.com/tools.nanobox.io/boxes/virtualbox/nanobox-boot2docker.md5"); err != nil {
		Config.Fatal("[commands/box/update] Util.MD5sMatch() failed - ", err.Error())
	}

	// if the local md5 does not match the remote md5 it's either wrong or old;
	// either way download the newest nanobox-boot2docker
	if !match {
		fmt.Printf(stylish.Bullet("Updating nanobox image..."))

		// update nanobox-boot2docker
		if err := Vagrant.Update(); err != nil {
			Config.Fatal("[commands/box/update] Vagrant.Update() failed - ", err.Error())
		}
	}

	// ensure the newly downloaded nanobox-boot2docker.box matches the remote one
	if match, err = Util.MD5sMatch(Config.Root()+"/nanobox-boot2docker.box", "https://s3.amazonaws.com/tools.nanobox.io/boxes/virtualbox/nanobox-boot2docker.md5"); err != nil {
		Config.Fatal("[commands/box/update] Util.MD5sMatch() failed - ", err.Error())
	}

	// if it doesn't match this time it's the wrong one
	if !match {
		fmt.Println("MD5 checksum failed! Your nanobox-boot2docker is not ours!")
	}
}
