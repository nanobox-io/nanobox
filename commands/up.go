// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"github.com/spf13/cobra"
)

//
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Resumes the halted/suspended nanobox VM",
	Long: `
Description:
  Runs 'nanobox create' and then 'nanobox deploy'`,

	Run: nanoUp,
}

//
func init() {
	upCmd.Flags().BoolVarP(&fWatch, "watch", "w", false, "Watches your app for file changes")
}

// nanoUp runs 'vagrant up'
func nanoUp(ccmd *cobra.Command, args []string) {

	// run a create command to create a Vagrantfile and boot the VM...
	nanoCreate(nil, args)

	// upgrade all nanobox docker images
	imagesUpdate(nil, args)

	// ...issue a deploy...
	nanoDeploy(nil, args)

	// ...begin watching the file system for changes
	if fWatch {
		nanoWatch(nil, args)
	}
}
