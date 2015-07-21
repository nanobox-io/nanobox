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

// ResumeCommand satisfies the Command interface
type ResumeCommand struct{}

// Help prints detailed help text for the app list command
func (c *ResumeCommand) Help() {
	ui.CPrint(`
Description:
  Resumes a halted/suspened nanobox VM

Usage:
  nanobox resume
  `)
}

// Run resumes the specified virtual machine
func (c *ResumeCommand) Run(opts []string) {

	// run 'vagrant resume'
	fmt.Printf(stylish.ProcessStart("resuming nanobox vm"))
	runVagrantCommand(exec.Command("vagrant", "resume"))
	fmt.Printf(stylish.ProcessEnd())
}
