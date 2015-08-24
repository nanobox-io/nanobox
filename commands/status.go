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

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Displays all current nanobox VM's",
	Long: `
Description:
  Displays all current nanobox VM's`,

	Run: nanoStatus,
}

// nanoStatus runs 'vagrant status'
func nanoStatus(ccmd *cobra.Command, args []string) {

	// run an init to ensure there is a Vagrantfile
	nanoInit(nil, args)

	fmt.Printf(stylish.ProcessStart("requesting nanobox vms"))
	if err := runVagrantCommand(exec.Command("vagrant", "status")); err != nil {
		ui.LogFatal("[commands/status] runVagrantCommand() failed", err)
	}
	fmt.Printf(stylish.ProcessEnd())
}
