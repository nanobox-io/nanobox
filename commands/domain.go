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
	"github.com/pagodabox/nanobox-golang-stylish"
)

// DomainCommand satisfies the Command interface
type DomainCommand struct{}

// Help prints detailed help text for the app list command
func (c *DomainCommand) Help() {
	ui.CPrint("Intended for internal use only")
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
		ui.LogFatal("[commands.domain] flags.Parse() failed", err)
	}

	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
	defer f.Close()

	//
	if err != nil {
		ui.LogFatal("[commands.domain] os.OpenFile() failed", err)
	}

	//
	switch {

	// write the entry to the hosts file
	case fWrite:
		entry := fmt.Sprintf("\n%-15v   %s.%s # '%v' private network (added by nanobox)", config.Boxfile.IP, config.App, config.Boxfile.Domain, config.App)

		if _, err := f.WriteString(entry); err != nil {
			ui.LogFatal("[commands.domain] WriteString() failed", err)
		}

		fmt.Println(stylish.Bullet(config.App + ".nano.dev added to /etc/hosts"))

	//
	case fRemove:

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

		// write back the contents of the hosts file minus the removed entry
		if err := ioutil.WriteFile("/etc/hosts", []byte(contents), 0644); err != nil {
			ui.LogFatal("[commands.destroy] ioutil.WriteFile failed", err)
		}

		fmt.Println(stylish.Bullet(config.App + ".nano.dev removed from /etc/hosts"))
	}
}

// modifyHosts
func modifyHosts(mod string) {

	// attempt to open /etc/hosts file...
	f, err := os.OpenFile("/etc/hosts", os.O_RDWR|os.O_APPEND, 0644)
	defer f.Close()

	//
	if err != nil {

		// if nanobox doesn't have permission to modify the hosts file, it needs to
		// request it
		if perm := os.IsPermission(err); perm == true {

			entry := fmt.Sprintf("%-15v   %s # '%v' private network (added by nanobox)", config.Boxfile.IP, config.Boxfile.Domain, config.App)

			//
			switch mod {

			// add
			case "w":
				fmt.Printf(`
Nanobox needs your permission to add the following entry to your /etc/hosts file:
%v

`, entry)

			// delete
			case "x":
				fmt.Printf(`
Nanobox needs your permission to remove the following entry from your /etc/hosts file:
%v

`, entry)
			}

			//
			cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v domain -%v", os.Args[0], mod))

			// connect standard in/outputs
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			// run command
			if err := cmd.Run(); err != nil {
				ui.LogFatal("[commands.domain] cmd.Run() failed", err)
			}

			//
		} else {
			ui.LogFatal("[commands.domain] os.OpenFile() failed", err)
		}
	}
}
