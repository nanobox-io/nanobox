// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"flag"

	"github.com/pagodabox/nanobox-cli/ui"
)

// UpCommand satisfies the Command interface
type UpCommand struct{}

// Help
func (c *UpCommand) Help() {
	ui.CPrint(`
Description:
  Runs 'nanobox create' and then 'nanobox deploy'

Usage:
  nanobox up

Options:
	-w, --watch
		Watches your app for file changes
  `)
}

// Run creates the specified virtual machine and issues a deploy to it
func (c *UpCommand) Run(opts []string) {

	// flags
	flags := flag.NewFlagSet("flags", flag.ContinueOnError)
	flags.Usage = func() { c.Help() }

	var fWatch bool
	flags.BoolVar(&fWatch, "w", false, "")
	flags.BoolVar(&fWatch, "watch", false, "")

	if err := flags.Parse(opts); err != nil {
		ui.LogFatal("[commands.destroy] flags.Parse() failed", err)
	}

	// run a create command to create a Vagrantfile and boot the VM...
	create := CreateCommand{}
	create.Run(opts)

	// ...issue a deploy...
	deploy := DeployCommand{}
	deploy.Run(opts)

	// ...begin watching the file system for changes
	if fWatch {
		watch := WatchCommand{}
		watch.Run(opts)
	}
}
