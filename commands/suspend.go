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

	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var suspendCmd = &cobra.Command{
	Use:   "suspend",
	Short: "Suspends the nanobox VM",
	Long: `
Description:
  Suspends the nanobox VM by issuing a "vagrant suspend"`,

	Run: nanoSuspend,
}

// nanoSuspend runs 'vagrant suspend'
func nanoSuspend(ccmd *cobra.Command, args []string) {

	// run an init to ensure there is a Vagrantfile
	nanoInit(nil, args)

	fmt.Printf(stylish.ProcessStart("suspending nanobox vm"))
	if err := util.RunVagrantCommand(exec.Command("vagrant", "suspend")); err != nil {
		util.Fatal("[commands/suspend] util.RunVagrantCommand() failed", err)
	}
	fmt.Printf(stylish.ProcessEnd())
}
