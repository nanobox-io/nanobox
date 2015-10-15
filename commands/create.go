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
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/file/hosts"
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var createCmd = &cobra.Command{
	Hidden: true,

	Use:   "create",
	Short: "Creates a new nanobox",
	Long:  ``,

	PreRun: initialize,
	Run:    create,
}

//
func init() {
	createCmd.Flags().BoolVarP(&fAddEntry, "add-entry", "", false, "")
	createCmd.Flags().MarkHidden("add-entry")
}

//
// create
func create(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	// if the command is being run with the "add" flag, it means an entry needs to
	// be added to the hosts file and execution yielded back to the parent
	if fAddEntry {
		hosts.AddDomain()
		os.Exit(0) // this exits the sudoed (child) created, not the parent proccess
	}

	// boot the vm
	fmt.Printf(stylish.Bullet("Creating a nanobox"))
	if err := vagrant.Up(); err != nil {
		vagrant.Fatal("[commands/create] vagrant.Up() failed - ", err.Error())
	}

	// after the machine boots, update the docker images
	updateImages(nil, args)

	// add the entry if needed
	if !hosts.HasDomain() {
		sudo("create --add-entry", fmt.Sprintf("Adding %v domain to hosts file", config.Nanofile.Domain))
	}

	// if devmode is detected, the machine needs to be rebooted for devmode to take
	// effect
	if config.Devmode {
		fmt.Printf(stylish.Bullet("Rebooting machine to finalize 'devmode' configuration..."))
		reload(nil, args)
	}
}
