// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"os/exec"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// SuspendCommand satisfies the Command interface
type SuspendCommand struct{}

// Help prints detailed help text for the app list command
func (c *SuspendCommand) Help() {
	ui.CPrint(`
Description:
  Suspends the current nanobox VM

Usage:
  nanobox suspend
  `)
}

// Run suspends the specified virtual machines
func (c *SuspendCommand) Run(opts []string) {

	// run 'vagrant suspend'
	fmt.Printf(stylish.ProcessStart("suspending nanobox vm"))
	runVagrantCommand(exec.Command("vagrant", "suspend"))
	fmt.Printf(stylish.ProcessEnd())
}
