// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"
	// "fmt"
	"os/exec"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/ui"
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
		config.Console.Warn("[commands.halt.Run] Issuing force delete.")
		cmd = exec.Command("vagrant", "halt", "--force")

		//
	} else {
		config.Console.Warn("[commands.halt.Run] Issuing confirm delete.")
		cmd = exec.Command("vagrant", "halt")
	}

	// run command
	if err := runVagrantCommand(cmd); err != nil {
		ui.LogFatal("[commands.halt] runVagrantCommand() failed", err)
	}

}
