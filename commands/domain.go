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

// DomainCommand satisfies the Command interface
type DomainCommand struct{}

// Help prints detailed help text for the app list command
func (c *DomainCommand) Help() {
	ui.CPrint(`
Description:
  Runs a specific sudo command

Usage:
  Intended for internal use only

Options:
  -w, --write
    Write the nanobox private network domain to the hosts file

  -x, --remove
    Remove the nanobox private network domain from the hosts file
  `)
}

// Run halts the specified virtual machine
func (c *DomainCommand) Run(opts []string) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fWrite bool
	flags.BoolVar(&fWrite, "w", false, "")
	flags.BoolVar(&fWrite, "write", false, "")

	var fRemove bool
	flags.BoolVar(&fRemove, "x", false, "")
	flags.BoolVar(&fRemove, "remove", false, "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.halt] flags.Parse() failed", err)
	}

	// attempt to open /etc/hosts file...
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)

	// if this command is run manually nanobox-cli may not have permission to
	// manipulate the hosts file, so we'll have to do what we do in create/destroy
	// and re-run the command as sudo
	if err != nil {

		//
		if perm := os.IsPermission(err); perm == true {
			//
			cmd := exec.Command("/bin/sh", "-c", "sudo "+os.Args[0])

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

	defer f.Close()

	//
	if fWrite {
		config.Console.Info("Adding '%v' private network to hosts file...", config.App)

		//
		entry := fmt.Sprintf("%-15v   %s.%s # '%v' private network (added by nanobox)", config.Boxfile.IP, config.App, config.Boxfile.Domain, config.App)

		// write the entry to the hosts file
		if _, err := f.WriteString(entry); err != nil {
			ui.LogFatal("[commands.sudo] WriteString() failed", err)
		}

		//
	} else if fRemove {
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

		//
	} else {
		config.Console.Fatal("Missing flag, please re-run the command and provide a valid flag.")
		os.Exit(1)
	}

}
