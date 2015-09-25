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
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Suspends the nanobox VM",
	Long: `
Description:
  Suspends the nanobox VM by issuing a "vagrant suspend"`,

	PreRun: projectIsCreated,
	Run:    nanoDown,
}

// nanoDown runs 'vagrant suspend'
func nanoDown(ccmd *cobra.Command, args []string) {

	// run an init to ensure there is a Vagrantfile
	nanoInit(nil, args)

	// if the CLI is running in background mode dont suspend the VM
	if fBackground {
		os.Exit(0)
	}

	//
	fmt.Printf(stylish.Bullet("Suspending nanobox VM..."))
	if err := util.RunVagrantCommand(exec.Command("vagrant", "suspend")); err != nil {
		config.Fatal("[commands/suspend] util.RunVagrantCommand() failed", err.Error())
	}
	fmt.Printf(stylish.Bullet("Exiting"))
}
