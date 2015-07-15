// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
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

	//
	// if force is not passed, confirm the decision to delete...
	if !fForce {
		fmt.Printf("------------------------- !! DANGER ZONE !! -------------------------\n\n")

		// prompt for confirmation...
		switch ui.Prompt("Are you sure you want to delete this VM (y/N)? ") {

		// if positive confirmation, proceed and destroy
		case "Y", "y":
			fmt.Printf(stylish.Bullet("Delete confirmed, continuing..."))
			destroy()

		// if negative confirmation, exit w/o destroying
		default:
			fmt.Printf(stylish.Bullet("Negative confirmation, app will not be deleted, exiting..."))
			os.Exit(0)
		}
	}

	// if force is passed, destroy...
	fmt.Printf(stylish.Bullet("Force delete detected, continuing..."))
	destroy()
}

// destroy
func destroy() {

	//
	// attempt to remove the associated entry, regardless of if it's there or not
	modifyHosts("x")

	//
	// destroy the vm...
	fmt.Printf(stylish.ProcessStart("destroying nanobox vm"))
	cmd := exec.Command("vagrant", "destroy", "--force")
	runVagrantCommand(cmd)

	// remove app; this needs to happen after the Vagrant command so that the app
	// isn't just created again
	fmt.Printf(stylish.Bullet("Deleting all nanobox files at: " + config.AppDir))
	if err := os.RemoveAll(config.AppDir); err != nil {
		ui.LogFatal("[commands.destroy] os.RemoveAll() failed", err)
	}
	fmt.Printf(stylish.ProcessEnd())
}
