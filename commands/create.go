// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// CreateCommand satisfies the Command interface
type CreateCommand struct{}

// Help prints detailed help text for the app list command
func (c *CreateCommand) Help() {
	ui.CPrint(`
Description:
  Runs an 'init' and starts a nanobox VM

Usage:
  nanobox create
  `)
}

// Run creates the specified virtual machine
func (c *CreateCommand) Run(opts []string) {

	// run an init to create a Vagrantfile...
	// init := InitCommand{}
	// init.Run(opts)

	//
	// open the /etc/hosts file for scanning...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		ui.LogFatal("[commands.create] os.Open() failed", err)
	}
	defer f.Close()

	// a new scanner for scanning the /etc/hosts file
	scanner := bufio.NewScanner(f)

	// determines whether or not an entry needs to be added to the /etc/hosts file
	// (an entry will be added unless it's confirmed that it's not needed)
	addEntry := true

	// scan hosts file looking for an entry corresponding to this app...
	for scanner.Scan() {

		// if an entry with the IP is detected, flag the entry as not needed
		if strings.HasPrefix(scanner.Text(), config.Boxfile.IP) {
			addEntry = false
		}
	}

	// add the entry as needed
	if addEntry {
		modifyHosts("w")
	}

	//
	// boot the vm
	fmt.Printf(stylish.ProcessStart("starting nanobox vm"))

	cmd := exec.Command("vagrant", "up")
	runVagrantCommand(cmd)

	fmt.Printf(stylish.ProcessEnd("nanobox vm running"))

}
