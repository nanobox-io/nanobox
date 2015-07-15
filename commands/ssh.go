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

// SSHCommand satisfies the Command interface
type SSHCommand struct{}

// Help prints detailed help text for the app list command
func (c *SSHCommand) Help() {
	ui.CPrint(`
Description:
  SSH into the virtual machine

Usage:
  nanobox ssh
  `)
}

// Run sshes into the virtual machine
func (c *SSHCommand) Run(opts []string) {

	// run 'vagrant ssh'
	fmt.Printf(stylish.ProcessStart("sshing into vm"))
	cmd := exec.Command("vagrant", "ssh")
	runVagrantCommand(cmd)
	fmt.Printf(stylish.ProcessEnd())
}
