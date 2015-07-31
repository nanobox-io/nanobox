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

// ReloadCommand satisfies the Command interface
type ReloadCommand struct{}

// Help
func (c *ReloadCommand) Help() {
	ui.CPrint(`
Description:
  Reloads the nanobox VM by issuing a "vagrant reload --provision"

Usage:
  nanobox reload
  `)
}

// Run reloads the specified virtual machine
func (c *ReloadCommand) Run(opts []string) {

	// run 'vagrant reload --provision'
	fmt.Printf(stylish.ProcessStart("reloading nanobox vm"))
	runVagrantCommand(exec.Command("vagrant", "reload", "--provision"))
	fmt.Printf(stylish.ProcessEnd())
}
