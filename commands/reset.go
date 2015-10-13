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
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var resetCmd = &cobra.Command{
	Hidden: true,

	Use:   "reset",
	Short: "Reloads the nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    reset,
}

// reset runs 'vagrant destroy', 'vagrant up', and a deploy
func reset(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	fmt.Printf(stylish.Bullet("Resetting nanobox..."))

	// destroy vm
	if err := vagrant.Destroy(); err != nil {
		config.Fatal("[commands/reload] failed - ", err.Error())
	}

	// create vm
	if err := vagrant.Up(); err != nil {
		config.Fatal("[commands/reload] failed - ", err.Error())
	}

	// issue a deploy
	deploy(nil, args)
}
