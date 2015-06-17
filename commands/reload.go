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

// ReloadCommand satisfies the Command interface
type ReloadCommand struct{}

// Help prints detailed help text for the app list command
func (c *ReloadCommand) Help() {
	ui.CPrint(`
Description:
  Reloads the Nanobox virtual machine

Usage:
  nanobox resume

Options:
  -p, --provision
    Also runs provisioners
  `)
}

// Run resumes the specified virtual machine
func (c *ReloadCommand) Run(opts []string) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fProvision bool
	flags.BoolVar(&fProvision, "p", false, "")
	flags.BoolVar(&fProvision, "provision", false, "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.reload] flags.Parse() failed", err)
	}

	cmd := &exec.Cmd{}

	// vagrant reload
	if fProvision {
		config.Console.Warn("[commands.reload.Run] Reloading and provisioning...")
		cmd = exec.Command("vagrant", "reload", "--provision")

		//
	} else {
		config.Console.Warn("[commands.reload.Run] Reloading...")
		cmd = exec.Command("vagrant", "reload")
	}

	// run command
	if err := runVagrantCommand(cmd); err != nil {
		ui.LogFatal("[commands.reload] runVagrantCommand() failed", err)
	}

}
