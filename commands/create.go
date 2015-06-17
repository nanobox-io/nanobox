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

// Run destroys the specified virtual machine
func (c *CreateCommand) Run(opts []string) {

	// run an init to create a Vagrantfile...
	init := InitCommand{}
	init.Run(opts)

	// boot the machine
	// NOTE: I want more time to think about the order of things here. Ideally I
	// want to break the writing portion of this out into it's own thing, but I
	// haven't quite come up with the best solution (that I like)
	cmd := exec.Command("vagrant", "up")

	if err := runVagrantCommand(cmd); err != nil {
		ui.LogFatal("[commands.create] runVagrantCommand() failed", err)
	}

	// assume we'll need to add an entry...
	addEntry := true

	// ...then check if we actually need to...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		ui.LogFatal("[commands.create] os.Open() failed", err)
	}

	defer f.Close()

	// ...read hosts file looking for an entry corresponding to this app...
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// ...don't write to the file if an entry with the IP is detected
		if strings.HasPrefix(scanner.Text(), config.Boxfile.IP) {
			addEntry = false
		}
	}

	// entry needed...
	if addEntry {

		//
		entry := fmt.Sprintf("\n%v %v # '%v' private network (added by nanobox)", config.Boxfile.IP, config.Boxfile.Domain, config.App)

		// attempt to open /etc/hosts file...
		f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)

		// ...if we're unable to open, we'll assume it's because we don't have permission
		if err != nil {

			fmt.Printf(`
Nanobox needs your permission to write the following entry into your /etc/hosts file:
	%v

	`, entry)

			// re-run the command as sudo so we can write to /etc/hosts
			if err := sudo(); err != nil {
				ui.LogFatal("[commands.create] sudo() failed", err)
			}

			os.Exit(0)
		}

		defer f.Close()

		config.Console.Info("Adding '%v' private network to hosts file...", config.App)

		// write the entry to the hosts file
		if _, err := f.WriteString(entry); err != nil {
			ui.LogFatal("[commands.create] WriteString() failed", err)
		}
	}

}
