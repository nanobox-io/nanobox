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
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroys the nanobox VM",
	Long: `
Description:
  Destroys the nanobox VM by issuing a "vagrant destroy"`,

	Run: nanoDestroy,
}

// nanoDestroy
func nanoDestroy(ccmd *cobra.Command, args []string) {

	// if nanobox is unable to access the /etc/hosts file it will need to re-run
	// the command as sudo
	if util.AccessDenied() {
		util.SudoExec("destroy", "Attempting to remove nano.dev domain from hosts file")
		os.Exit(0)
	}

	//
	// if force is not passed, confirm the decision to delete...
	if !fForce {
		fmt.Printf("------------------------- !! DANGER ZONE !! -------------------------\n\n")

		// prompt for confirmation...
		switch util.Prompt("Are you sure you want to delete this VM (y/N)? ") {

		// if positive confirmation, proceed and destroy
		case "Y", "y":
			fmt.Printf(stylish.Bullet("Delete confirmed, continuing..."))

		// if negative confirmation, exit w/o destroying
		default:
			os.Exit(0)
		}
	}

	// attempt to remove the associated entry, regardless of if it's there or not
	util.RemoveDevDomain()

	//
	// destroy the vm; this needs to happen first to ensure there is a Vagrantfile
	// to run the command with
	fmt.Printf(stylish.ProcessStart("destroying nanobox vm"))
	if err := runVagrantCommand(exec.Command("vagrant", "destroy", "--force")); err != nil {
		if err == err.(*os.PathError) {
			return
		} else {
			util.LogFatal("[commands/destroy] runVagrantCommand() failed", err)
		}
	}

	// remove app; this needs to happen last so that the app isn't just created
	// again while running the vagrant command
	fmt.Printf(stylish.Bullet("Deleting all nanobox files at: " + config.AppDir))
	if err := os.RemoveAll(config.AppDir); err != nil {
		util.LogFatal("[commands/destroy] os.RemoveAll() failed", err)
	}
	fmt.Printf(stylish.ProcessEnd())
}
