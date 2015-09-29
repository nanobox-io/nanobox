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

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
)

//
var nanoboxDevCmd = &cobra.Command{
	Use:   "dev",
	Short: "Starts the nanobox VM, provisions app, & opens an interactive terminal",
	Long:  ``,

	PreRun:  bootVM,
	Run:     nanoboxDev,
	PostRun: saveVM,
}

//
func init() {
	nanoboxDevCmd.Flags().BoolVarP(&fRebuild, "rebuild", "", false, "Rebuilds")
}

// nanoboxDev
func nanoboxDev(ccmd *cobra.Command, args []string) {

	// PreRun: bootVM

	//
	switch {

	// if the vm is 'new' deploy and update images
	case util.VagrantStatus() == "not created" || !config.VMfile.HasDeployed():
		imagesUpdate(nil, args)
		nanoDeploy(nil, args)

	// if fRebuild is detected only deploy
	case fRebuild:
		nanoDeploy(nil, args)
	}

	//
	nanoboxConsole(nil, args)

	// PostRun: saveVM
}
