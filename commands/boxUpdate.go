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
var boxUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,

	Run: boxUpdate,
}

// boxUpdate
func boxUpdate(ccmd *cobra.Command, args []string) {

	// check to make sure there is a box
	boxInstall(nil, args)

	// if the local md5 does not match the remote md5...
	if util.MD5sMatch(config.Root+"/nanobox-boot2docker.md5", "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5") {
		return
	}

	// ...download the new image
	fmt.Printf(stylish.Bullet("Updating nanobox image..."))
	util.VMDownload()
}
