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

// StatusCommand satisfies the Command interface
type StatusCommand struct{}

// Help prints detailed help text for the app list command
func (c *StatusCommand) Help() {
	ui.CPrint(`
Description:
  Display all current nanobox VM's

Usage:
  nanobox status
  `)
}

// Run display status of all virtual machines
func (c *StatusCommand) Run(opts []string) {

	// run 'vagrant status'
	fmt.Printf(stylish.ProcessStart("requesting nanobox vms"))
	runVagrantCommand(exec.Command("vagrant", "status"))
	fmt.Printf(stylish.ProcessEnd())
}
