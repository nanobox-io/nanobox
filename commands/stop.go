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
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
)

//
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Suspends the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    stop,
}

// stop runs 'vagrant suspend'
func stop(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	//
	if err := vagrant.Suspend(); err != nil {
		vagrant.Fatal("[commands/stop] vagrant.Suspend() failed - ", err.Error())
	}

	// set the mode to be forground next time the machine boots
	config.VMfile.ModeIs("foreground")
}
