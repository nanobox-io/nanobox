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
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Suspends the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    down,
}

// down runs 'vagrant suspend'
func down(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	//
	if err := vagrant.Suspend(); err != nil {
		config.Fatal("[commands/nanoboxDown] failed - ", err.Error())
	}

	// set the mode to be forground next time the machine boots
	config.VMfile.ModeIs("foreground")
}
