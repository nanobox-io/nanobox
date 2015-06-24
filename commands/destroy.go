// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
)

// DestroyCommand satisfies the Command interface
type DestroyCommand struct{}

// Help prints detailed help text for the app list command
func (c *DestroyCommand) Help() {
	ui.CPrint(`
Description:
  Destroys the current nanobox VM

Usage:
  nanobox destroy

Options:
  -f, --force
    A force destroy [red]skips confirmation... use responsibly[reset]!
  `)
}

// Run destroys the specified virtual machine
func (c *DestroyCommand) Run(opts []string) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fForce bool
	flags.BoolVar(&fForce, "f", false, "")
	flags.BoolVar(&fForce, "force", false, "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	// assume we wont need to remove an entry...
	removeEntry := false

	// ...then check if we actually need to...
	f, err := os.Open("/etc/hosts")
	if err != nil {
		ui.LogFatal("[commands.destroy] os.Open() failed", err)
	}

	defer f.Close()

	// ...read hosts file looking for an entry corresponding to this app...
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		// ...remove from the file if an entry with the IP is detected
		if strings.HasPrefix(scanner.Text(), config.Boxfile.IP) {
			removeEntry = true
		}
	}

	// remove entry
	if removeEntry {

		// attempt to open /etc/hosts file...
		f, err := os.OpenFile("/etc/hosts", os.O_RDWR, 0644)
		defer f.Close()

		// ...if we're unable to open, we'll assume it's because we don't have permission
		if err != nil {

			//
			if perm := os.IsPermission(err); perm == true {

				//
				cmd := exec.Command("/bin/sh", "-c", "sudo "+os.Args[0]+" domain -x")

				// connect standard in/outputs
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				//
				fmt.Printf("\nNanobox needs your permission to remove the '%v' network from your /etc/hosts file\n", config.App)

				// run command
				if err := cmd.Run(); err != nil {
					ui.LogFatal("[commands.create] cmd.Run() failed", err)
				}

				//
			} else {
				ui.LogFatal("[commands.domain] os.OpenFile() failed", err)
			}
		}
	}

	cmd := &exec.Cmd{}

	// vagrant destroy
	if fForce {
		config.Console.Warn("[commands.destroy.Run] Issuing force delete.")
		cmd = exec.Command("vagrant", "destroy", "--force")

		//
	} else {
		config.Console.Warn("[commands.destroy.Run] Issuing confirm delete.")
		cmd = exec.Command("vagrant", "destroy")
	}

	// destroy the vm...
	if err := runVagrantCommand(cmd); err != nil {
		ui.LogFatal("[commands.destroy] runVagrantCommand() failed", err)
	}

	// destroy the project folder in /.nanobox with the Vagrantfile and .vagrant folder
	if err := os.RemoveAll(config.AppDir); err != nil {
		ui.LogFatal("[commands.destroy] os.RemoveAll() failed", err)
	}
}
