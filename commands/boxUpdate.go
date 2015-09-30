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

	// "github.com/nanobox-io/nanobox-cli/config"
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

	// check to make sure there is a box already
	boxInstall(nil, args)

	// if the local md5 doesn't match remote md5, download the box
	if util.VMLocalMD5() != util.VMRemoteMD5() {
		fmt.Printf(stylish.Bullet("Updating virtual machine"))
		util.VMDownload()
	}
}
