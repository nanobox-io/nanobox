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
// var sandboxCmd = &cobra.Command{
// 	Use:   "sanbox",
// 	Short: "",
// 	Long:  ``,
//
// 	Run: nanoSandbox,
// }

// nanoSandbox runs 'vagrant up'
func nanoSandbox(ccmd *cobra.Command, args []string) {

	//
	// run a create command to create a Vagrantfile and boot the VM...
	nanoCreate(nil, args)

	// upgrade all nanobox docker images
	imagesUpdate(nil, args)

	// ...issue a deploy...
	// fSandbox = true;
	nanoDeploy(nil, []string{"--sandbox"})
}
