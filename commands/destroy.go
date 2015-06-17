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
	"io/ioutil"
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

		// ...if we're unable to open, we'll assume it's because we don't have permission
		if err != nil {

			fmt.Printf("Nanobox needs your permission to remove the '%v' network from your /etc/hosts file\n", config.App)

			// re-run the command as sudo so we can write to /etc/hosts
			if err := sudo(); err != nil {
				ui.LogFatal("[commands.destroy] sudo() failed", err)
			}

			os.Exit(0)
		}

		defer f.Close()

		config.Console.Info("Removing '%v' private network to hosts file...", config.App)

		scanner := bufio.NewScanner(f)
		contents := ""

		// remove entry from /etc/hosts
		for scanner.Scan() {

			// if the line doesn't contain the entry add it back to what is going to be
			// re-written to the file
			if !strings.HasPrefix(scanner.Text(), config.Boxfile.IP) {
				contents += fmt.Sprintf("%s\n", scanner.Text())
			}
		}

		// write back the entirety of the hosts file minus the removed entry
		err = ioutil.WriteFile("/etc/hosts", []byte(contents), 0644)
		if err != nil {
			ui.LogFatal("[commands.destroy] ioutil.WriteFile failed", err)
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
