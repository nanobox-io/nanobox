// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

//
import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/box"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

//
var initCmd = &cobra.Command{
	Hidden: true,

	Use:   "init",
	Short: "Creates a nanobox-flavored Vagrantfile",
	Long:  ``,

	Run: initialize,
}

// initialize
func initialize(ccmd *cobra.Command, args []string) {

	// check to see if a box needs to be installed
	box.Install(nil, args)

	// creates a project folder at ~/.nanobox/apps/<name> (if it doesn't already
	// exists) where the Vagrantfile and .vagrant dir will live for each app
	if _, err := os.Stat(config.AppDir); err != nil {
		if err := os.Mkdir(config.AppDir, 0755); err != nil {
			config.Fatal("[commands/init] os.Mkdir() failed", err.Error())
		}
	}

	// 'parse' the .vmfile (either creating one, or parsing it)
	config.VMfile = config.ParseVMfile()

	//
	// generate a Vagrantfile at ~/.nanobox/apps/<app-name>/Vagrantfile
	// only if one doesn't already exist (unless forced)
	if !config.Force {
		if _, err := os.Stat(config.AppDir + "/Vagrantfile"); err == nil {
			return
		}
	}

	vagrant.Init()
}
