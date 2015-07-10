// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	"os/exec"

	"github.com/pagodabox/nanobox-cli/ui"
	"github.com/pagodabox/nanobox-golang-stylish"
)

// HaltCommand satisfies the Command interface
type HaltCommand struct{}

// Help prints detailed help text for the app list command
func (c *HaltCommand) Help() {
	ui.CPrint(`
Description:
  Halts the current nanobox VM

Usage:
  nanobox halt

Options:
  -f, --force
    A forced halt
  `)
}

// Run halts the specified virtual machine
func (c *HaltCommand) Run(opts []string) {
	stylish.ProcessStart("halting nanobox vm")

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fForce bool
	flags.BoolVar(&fForce, "f", false, "")
	flags.BoolVar(&fForce, "force", false, "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.halt] flags.Parse() failed", err)
	}

	cmd := &exec.Cmd{}

	// vagrant halt
	if fForce {
		stylish.Bullet("Issuing forced halt...")
		cmd = exec.Command("vagrant", "halt", "--force")

		//
	} else {
		stylish.Bullet("Issuing confirm halt...")
		cmd = exec.Command("vagrant", "halt")
	}

	// run command
	runVagrantCommand(cmd)

	stylish.ProcessEnd("nanobox vm halted")
}
