// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// SSHCommand satisfies the Command interface
type SSHCommand struct{}

// Help
func (c *SSHCommand) Help() {
	ui.CPrint(`
Description:
  SSHes into the nanobox VM by issuing a "vagrant ssh"

Usage:
  nanobox ssh
  `)
}

// Run
func (c *SSHCommand) Run(opts []string) {
	fmt.Printf(stylish.Bullet("SSHing into nanobox VM..."))

	// run an init to ensure there is a Vagrantfile
	init := InitCommand{}
	init.Run(nil)

	// run the command from ~/.nanobox/apps/<this app>
	if err := os.Chdir(config.AppDir); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] os.Chdir() failed", err)
	}

	cmd := exec.Command("vagrant", "ssh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//
	fmt.Printf(stylish.Bullet(fmt.Sprintf("running '%v'", strings.Trim(fmt.Sprint(cmd.Args), "[]"))))

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Run(); err != nil {
		ui.LogFatal("[commands.runVagrantCommand] cmd.Start() failed", err)
	}
}
