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
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
	"github.com/nanobox-io/nanobox-golang-stylish"
)

//
var nanoboxDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Suspends the nanobox VM",
	Long: `
Description:
  Suspends the nanobox VM by issuing a "vagrant suspend"`,

	PreRun: nanoInit,
	Run:    nanoboxDown,
}

// nanoboxDown runs 'vagrant suspend'
func nanoboxDown(ccmd *cobra.Command, args []string) {

	// PreRun: nanoInit

	//
	fmt.Printf("\n%s", stylish.Bullet("Suspending nanobox VM..."))
	if err := util.VagrantRun(exec.Command("vagrant", "suspend")); err != nil {
		config.Fatal("[commands/suspend] util.VagrantRun() failed", err.Error())
	}
	fmt.Printf(stylish.Bullet("Exiting"))

	// set the mode to be forground next time the machine boots
	config.VMfile.ModeIs("foreground")
}
